package sidecar

import (
	"fmt"
	"reflect"
	"strings"

	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	httpConnectionManager "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tcpProxy "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"istio.io/istio/pilot/pkg/networking/util"
)

type LDS struct {
	Address     string      `json:"address"`
	Port        uint32      `json:"port"`
	Match       interface{} `json:"match"`
	Destination string      `json:"destination"`
}

func ClustersToLDS(configDump *ConfigDump) []LDS {
	lds := make([]LDS, 0)
	listeners, err := configDump.GetListeners()
	if err != nil {
		return lds
	}
	for _, l := range listeners {
		address := retrieveListenerAddress(l)
		port := l.GetAddress().GetSocketAddress().GetPortValue()
		matchs := retrieveListenerMatches(l)
		for _, match := range matchs {
			lds = append(lds, LDS{
				Address:     address,
				Port:        port,
				Match:       match.match,
				Destination: match.destination,
			})
		}
	}
	return lds
}

func retrieveListenerAddress(l *listener.Listener) string {
	sockAddr := l.Address.GetSocketAddress()
	if sockAddr != nil {
		return sockAddr.Address
	}

	pipe := l.Address.GetPipe()
	if pipe != nil {
		return pipe.Path
	}

	return ""
}

type filterChain struct {
	match       string
	destination string
}

var protocolDescriptions = map[string][]string{
	"App: HTTP TLS":         []string{"http/1.0", "http/1.1", "h2c", "istio-http/1.0", "istio-http/1.1", "istio-h2"},
	"App: Istio HTTP Plain": []string{"istio", "istio-http/1.0", "istio-http/1.1", "istio-h2"},
	"App: TCP TLS":          []string{"istio-peer-exchange", "istio"},
	"App: HTTP":             []string{"http/1.0", "http/1.1", "h2c"},
}

func retrieveListenerMatches(l *listener.Listener) []filterChain {
	fc := l.FilterChains
	if l.DefaultFilterChain != nil {
		fc = append(fc, l.DefaultFilterChain)
	}

	resp := make([]filterChain, 0, len(fc))
	for _, f := range fc {
		match := f.FilterChainMatch
		if match == nil {
			match = &listener.FilterChainMatch{}
		}

		descriptions := make([]string, 0)

		if len(match.ServerNames) > 0 {
			descriptions = append(descriptions, fmt.Sprintf("SNI: %s", strings.Join(match.ServerNames, ",")))
		}

		if len(match.TransportProtocol) > 0 {
			descriptions = append(descriptions, fmt.Sprintf("Trans: %s", match.TransportProtocol))
		}

		if len(match.ApplicationProtocols) > 0 {
			found := false
			for protocolDescription, protocols := range protocolDescriptions {
				if reflect.DeepEqual(match.ApplicationProtocols, protocols) {
					found = true
					descriptions = append(descriptions, protocolDescription)
					break
				}
			}
			if !found {
				descriptions = append(descriptions, fmt.Sprintf("App: %s", strings.Join(match.ApplicationProtocols, ",")))
			}
		}

		port := ""
		if match.DestinationPort != nil {
			port = fmt.Sprintf(":%d", match.DestinationPort.GetValue())
		}
		if len(match.PrefixRanges) > 0 {
			pf := make([]string, 0)
			for _, p := range match.PrefixRanges {
				pf = append(pf, fmt.Sprintf("%s/%d", p.AddressPrefix, p.GetPrefixLen().GetValue()))
			}
			descriptions = append(descriptions, fmt.Sprintf("Addr: %s%s", strings.Join(pf, ","), port))
		} else if port != "" {
			descriptions = append(descriptions, fmt.Sprintf("Addr: *%s", port))
		}
		if len(descriptions) == 0 {
			descriptions = []string{"ALL"}
		}
		resp = append(resp, filterChain{
			destination: getFilterType(f.GetFilters()),
			match:       strings.Join(descriptions, "; "),
		})
	}
	return resp
}

func getFilterType(filters []*listener.Filter) string {
	for _, filter := range filters {
		if filter.Name == wellknown.HTTPConnectionManager {
			httpProxy := &httpConnectionManager.HttpConnectionManager{}
			filter.GetTypedConfig().TypeUrl = "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager"
			err := filter.GetTypedConfig().UnmarshalTo(httpProxy)
			if err != nil {
				return err.Error()
			}
			if httpProxy.GetRouteConfig() != nil {
				return describeRouteConfig(httpProxy.GetRouteConfig())
			}
			if httpProxy.GetRds().GetRouteConfigName() != "" {
				return fmt.Sprintf("Route: %s", httpProxy.GetRds().GetRouteConfigName())
			}
			return "HTTP"
		} else if filter.Name == wellknown.TCPProxy {
			if !strings.Contains(string(filter.GetTypedConfig().GetValue()), util.BlackHoleCluster) {
				tp := &tcpProxy.TcpProxy{}
				// Allow Unmarshal to work even if Envoy and istioctl are different
				filter.GetTypedConfig().TypeUrl = "type.googleapis.com/envoy.extensions.filters.network.tcp_proxy.v3.TcpProxy"
				err := filter.GetTypedConfig().UnmarshalTo(tp)
				if err != nil {
					return err.Error()
				}
				if strings.Contains(tp.GetCluster(), "Cluster") {
					return tp.GetCluster()
				}
				return fmt.Sprintf("Cluster: %s", tp.GetCluster())
			}
		}
	}
	return "Non-HTTP/Non-TCP"
}

func describeRouteConfig(route *route.RouteConfiguration) string {
	if cluster := getMatchAllCluster(route); cluster != "" {
		return cluster
	}
	vhosts := make([]string, 0)
	for _, vh := range route.GetVirtualHosts() {
		if describeDomains(vh) == "" {
			vhosts = append(vhosts, describeRoutes(vh))
		} else {
			vhosts = append(vhosts, fmt.Sprintf("%s %s", describeDomains(vh), describeRoutes(vh)))
		}
	}
	return fmt.Sprintf("Inline Route: %s", strings.Join(vhosts, "; "))
}

func getMatchAllCluster(er *route.RouteConfiguration) string {
	if len(er.GetVirtualHosts()) != 1 {
		return ""
	}
	vh := er.GetVirtualHosts()[0]
	if !reflect.DeepEqual(vh.Domains, []string{"*"}) {
		return ""
	}
	if len(vh.GetRoutes()) != 1 {
		return ""
	}
	r := vh.GetRoutes()[0]
	if r.GetMatch().GetPrefix() != "/" {
		return ""
	}
	a, ok := r.GetAction().(*route.Route_Route)
	if !ok {
		return ""
	}
	cl, ok := a.Route.ClusterSpecifier.(*route.RouteAction_Cluster)
	if !ok {
		return ""
	}
	if strings.Contains(cl.Cluster, "Cluster") {
		return cl.Cluster
	}
	return fmt.Sprintf("Cluster: %s", cl.Cluster)
}

func describeDomains(vh *route.VirtualHost) string {
	if len(vh.GetDomains()) == 1 && vh.GetDomains()[0] == "*" {
		return ""
	}
	return strings.Join(vh.GetDomains(), "/")
}

func describeRoutes(vh *route.VirtualHost) string {
	routes := make([]string, 0, len(vh.GetRoutes()))
	for _, r := range vh.GetRoutes() {
		routes = append(routes, describeMatch(r.GetMatch()))
	}
	return strings.Join(routes, ", ")
}

func describeMatch(match *route.RouteMatch) string {
	conds := make([]string, 0)
	if match.GetPrefix() != "" {
		conds = append(conds, fmt.Sprintf("%s*", match.GetPrefix()))
	}
	if match.GetPath() != "" {
		conds = append(conds, match.GetPath())
	}
	if match.GetSafeRegex() != nil {
		conds = append(conds, fmt.Sprintf("regex %s", match.GetSafeRegex().Regex))
	}
	// Ignore headers
	return strings.Join(conds, " ")
}
