package sidecar

import (
	"encoding/json"

	adminapi "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	"istio.io/istio/pkg/util/protomarshal"
)

type Cluster struct {
	*adminapi.Clusters
}

func NewCluster(config []byte) (*Cluster, error) {
	cluster := &Cluster{}
	err := json.Unmarshal(config, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (s *Cluster) MarshalJSON() ([]byte, error) {
	return protomarshal.Marshal(s)
}

func (s *Cluster) UnmarshalJSON(b []byte) error {
	cd := &adminapi.Clusters{}
	err := protomarshal.UnmarshalAllowUnknown(b, cd)
	*s = Cluster{cd}
	return err
}
