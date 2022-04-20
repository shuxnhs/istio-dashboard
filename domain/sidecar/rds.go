package sidecar

import (
	"fmt"
	"strings"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"istio.io/istio/pilot/pkg/networking/util"
	"istio.io/istio/pkg/util/sets"
)

type RDS struct {
	Name           string   `json:"name"`
	Domains        []string `json:"domains"`
	Match          string   `json:"match"`
	VirtualService string   `json:"virtualService"`
}

func ClustersToRDS(configDump *ConfigDump) []RDS {
	rds := make([]RDS, 0)
	routes, err := configDump.GetRouters()
	if err != nil {
		return rds
	}
	for _, r := range routes {
		name := r.GetName()
		for _, hosts := range r.VirtualHosts {
			if len(hosts.Routes) == 0 {
				rds = append(rds, RDS{
					Name:           name,
					Domains:        describeRouteDomains(hosts.Domains),
					Match:          "/*",
					VirtualService: "404",
				})
			}
			for _, r := range hosts.Routes {
				if !isPassThrough(r.GetAction()) {
					rds = append(rds, RDS{
						Name:           name,
						Domains:        describeRouteDomains(hosts.GetDomains()),
						Match:          describeMatch(r.GetMatch()),
						VirtualService: describeManagement(r.GetMetadata()),
					})
				}
			}

		}

	}
	return rds
}

func describeRouteDomains(domains []string) []string {
	withoutPort := make([]string, 0, len(domains))
	for _, d := range domains {
		if !strings.Contains(d, ":") {
			withoutPort = append(withoutPort, d)
			// if the domain contains IPv6, such as [fd00:10:96::7fc7] and [fd00:10:96::7fc7]:8090
		} else if strings.Count(d, ":") > 2 {
			// if the domain is only a IPv6 address, such as [fd00:10:96::7fc7], append it
			if strings.HasSuffix(d, "]") {
				withoutPort = append(withoutPort, d)
			}
		}
	}
	return unExpandDomains(withoutPort)
}

func isPassThrough(action interface{}) bool {
	a, ok := action.(*route.Route_Route)
	if !ok {
		return false
	}
	cl, ok := a.Route.ClusterSpecifier.(*route.RouteAction_Cluster)
	if !ok {
		return false
	}
	return cl.Cluster == "PassthroughCluster"
}

func describeManagement(metadata *core.Metadata) string {
	if metadata == nil {
		return ""
	}
	istioMetadata, ok := metadata.FilterMetadata[util.IstioMetadataKey]
	if !ok {
		return ""
	}
	config, ok := istioMetadata.Fields["config"]
	if !ok {
		return ""
	}
	return renderConfig(config.GetStringValue())
}

func renderConfig(configPath string) string {
	if strings.HasPrefix(configPath, "/apis/networking.istio.io/v1alpha3/namespaces/") {
		pieces := strings.Split(configPath, "/")
		if len(pieces) != 8 {
			return ""
		}
		return fmt.Sprintf("%s.%s", pieces[7], pieces[5])
	}
	return "<unknown>"
}

func unExpandDomains(domains []string) []string {
	unique := sets.New(domains...)
	shouldDelete := sets.New()
	for _, h := range domains {
		stripFull := strings.TrimSuffix(h, ".svc.cluster.local")
		if _, f := unique[stripFull]; f && stripFull != h {
			shouldDelete.Insert(h)
		}
		stripPartial := strings.TrimSuffix(h, ".svc")
		if _, f := unique[stripPartial]; f && stripPartial != h {
			shouldDelete.Insert(h)
		}
	}
	// Filter from original list to keep original order
	ret := make([]string, 0, len(domains))
	for _, h := range domains {
		if _, f := shouldDelete[h]; !f {
			ret = append(ret, h)
		}
	}
	return ret
}
