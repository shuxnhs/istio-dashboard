package sidecar

import (
	"encoding/json"
	"fmt"

	admin "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/golang/protobuf/ptypes/any"
	"istio.io/istio/pkg/util/protomarshal"
)

const (
	clusters  string = "type.googleapis.com/envoy.admin.v3.ClustersConfigDump"
	listeners string = "type.googleapis.com/envoy.admin.v3.ListenersConfigDump"
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
			c.Cluster.TypeUrl = "type.googleapis.com/envoy.config.cluster.v3.Cluster"
			if err := c.Cluster.UnmarshalTo(tmpCluster); err == nil {
				clusters = append(clusters, tmpCluster)
			}
		}
	}
	for _, c := range cd.DynamicActiveClusters {
		if c.GetCluster() != nil {
			tmpCluster := &cluster.Cluster{}
			c.Cluster.TypeUrl = "type.googleapis.com/envoy.config.cluster.v3.Cluster"
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
	cd, err := c.GetListenerConfigDump()
	if err != nil {
		return nil, err
	}
	for _, l := range cd.StaticListeners {
		tmpListener := &listener.Listener{}
		l.Listener.TypeUrl = "type.googleapis.com/envoy.config.listener.v3.Listener"
		if err := l.Listener.UnmarshalTo(tmpListener); err == nil {
			listeners = append(listeners, tmpListener)
		}
	}

	for _, l := range cd.DynamicListeners {
		tmpListener := &listener.Listener{}
		l.ActiveState.Listener.TypeUrl = "type.googleapis.com/envoy.config.listener.v3.Listener"
		if err := l.ActiveState.Listener.UnmarshalTo(tmpListener); err == nil {
			listeners = append(listeners, tmpListener)
		}
	}

	return listeners, nil
}
