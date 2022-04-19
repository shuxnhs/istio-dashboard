package sidecar

import (
	"encoding/json"
	"fmt"

	adminapi "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/golang/protobuf/ptypes/any"
	"istio.io/istio/pkg/util/protomarshal"
)

const clusters string = "type.googleapis.com/envoy.admin.v3.ClustersConfigDump"

type ConfigDump struct {
	*adminapi.ConfigDump
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
	cd := &adminapi.ConfigDump{}
	err := protomarshal.UnmarshalAllowUnknown(b, cd)
	*c = ConfigDump{cd}
	return err
}

func (c *ConfigDump) GetClusterConfigDump() (*adminapi.ClustersConfigDump, error) {
	var clusterDumpAny *any.Any
	for _, conf := range c.Configs {
		if conf.TypeUrl == clusters {
			clusterDumpAny = conf
		}
	}
	if clusterDumpAny == nil {
		return nil, fmt.Errorf("config dump has no configuration type %s", clusters)
	}

	clusterDump := &adminapi.ClustersConfigDump{}
	err := clusterDumpAny.UnmarshalTo(clusterDump)
	if err != nil {
		return nil, err
	}
	return clusterDump, nil
}

func (c *ConfigDump) GetCluster() ([]*cluster.Cluster, error) {
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
