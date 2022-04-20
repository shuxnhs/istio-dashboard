package sidecar

import (
	"encoding/json"
	"fmt"

	admin "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/ptypes/any"
	"istio.io/istio/pilot/pkg/xds/v3"
	"istio.io/istio/pkg/util/protomarshal"
)

const (
	clusters  string = "type.googleapis.com/envoy.admin.v3.ClustersConfigDump"
	listeners string = "type.googleapis.com/envoy.admin.v3.ListenersConfigDump"
	routes    string = "type.googleapis.com/envoy.admin.v3.RoutesConfigDump"
)

type ConfigDump struct {
	*admin.ConfigDump
}

func NewConfigDump(config []byte) (*ConfigDump, error) {
	configDump := &ConfigDump{}
	err := json.Unmarshal(config, configDump)
	if err != nil {
		return nil, err
	}
	return configDump, nil
}

func (c *ConfigDump) MarshalJSON() ([]byte, error) {
	return protomarshal.Marshal(c)
}

func (c *ConfigDump) UnmarshalJSON(b []byte) error {
	cd := &admin.ConfigDump{}
	err := protomarshal.UnmarshalAllowUnknown(b, cd)
	*c = ConfigDump{cd}
	return err
}

func (c *ConfigDump) GetClusterConfigDump() (*admin.ClustersConfigDump, error) {
	var clusterDumpAny *any.Any
	for _, conf := range c.Configs {
		if conf.TypeUrl == clusters {
			clusterDumpAny = conf
		}
	}
	if clusterDumpAny == nil {
		return nil, fmt.Errorf("config dump has no configuration type %s", clusters)
	}

	clusterDump := &admin.ClustersConfigDump{}
	err := clusterDumpAny.UnmarshalTo(clusterDump)
	if err != nil {
		return nil, err
	}
	return clusterDump, nil
}

func (c *ConfigDump) GetClusters() ([]*cluster.Cluster, error) {
	clusters := make([]*cluster.Cluster, 0)
	cd, err := c.GetClusterConfigDump()
	if err != nil {
		return nil, err
	}
	for _, c := range cd.StaticClusters {
		if c.GetCluster() != nil {
			tmpCluster := &cluster.Cluster{}
			c.Cluster.TypeUrl = v3.ClusterType
			if err := c.Cluster.UnmarshalTo(tmpCluster); err == nil {
				clusters = append(clusters, tmpCluster)
			}
		}
	}
	for _, c := range cd.DynamicActiveClusters {
		if c.GetCluster() != nil {
			tmpCluster := &cluster.Cluster{}
			c.Cluster.TypeUrl = v3.ClusterType
			if err := c.Cluster.UnmarshalTo(tmpCluster); err == nil {
				clusters = append(clusters, tmpCluster)
			}
		}
	}
	return clusters, nil
}

func (c *ConfigDump) GetListenerConfigDump() (*admin.ListenersConfigDump, error) {
	var listenerDumpAny *any.Any
	for _, conf := range c.Configs {
		if conf.TypeUrl == listeners {
			listenerDumpAny = conf
		}
	}
	if listenerDumpAny == nil {
		return nil, fmt.Errorf("config dump has no configuration type %s", clusters)
	}

	listenerDump := &admin.ListenersConfigDump{}
	err := listenerDumpAny.UnmarshalTo(listenerDump)
	if err != nil {
		return nil, err
	}
	return listenerDump, nil
}

func (c *ConfigDump) GetListeners() ([]*listener.Listener, error) {
	listeners := make([]*listener.Listener, 0)
	ld, err := c.GetListenerConfigDump()
	if err != nil {
		return nil, err
	}
	for _, l := range ld.StaticListeners {
		tmpListener := &listener.Listener{}
		l.Listener.TypeUrl = v3.ListenerType
		if err := l.Listener.UnmarshalTo(tmpListener); err == nil {
			listeners = append(listeners, tmpListener)
		}
	}

	for _, l := range ld.DynamicListeners {
		tmpListener := &listener.Listener{}
		l.ActiveState.Listener.TypeUrl = v3.ListenerType
		if err := l.ActiveState.Listener.UnmarshalTo(tmpListener); err == nil {
			listeners = append(listeners, tmpListener)
		}
	}

	return listeners, nil
}

func (c *ConfigDump) GetRouterConfigDump() (*admin.RoutesConfigDump, error) {
	var routeDumpAny *any.Any
	for _, conf := range c.Configs {
		if conf.TypeUrl == routes {
			routeDumpAny = conf
		}
	}
	if routeDumpAny == nil {
		return nil, fmt.Errorf("config dump has no configuration type %s", routes)
	}

	routeDump := &admin.RoutesConfigDump{}
	err := routeDumpAny.UnmarshalTo(routeDump)
	if err != nil {
		return nil, err
	}
	return routeDump, nil
}

func (c *ConfigDump) GetRouters() ([]*route.RouteConfiguration, error) {
	routes := make([]*route.RouteConfiguration, 0)
	rd, err := c.GetRouterConfigDump()
	if err != nil {
		return nil, err
	}
	for _, r := range rd.StaticRouteConfigs {
		tmpRoute := &route.RouteConfiguration{}
		r.RouteConfig.TypeUrl = v3.RouteType
		if err := r.RouteConfig.UnmarshalTo(tmpRoute); err == nil {
			routes = append(routes, tmpRoute)
		}
	}

	for _, r := range rd.DynamicRouteConfigs {
		tmpRoute := &route.RouteConfiguration{}
		r.RouteConfig.TypeUrl = v3.RouteType
		if err := r.RouteConfig.UnmarshalTo(tmpRoute); err == nil {
			routes = append(routes, tmpRoute)
		}
	}

	return routes, nil
}
