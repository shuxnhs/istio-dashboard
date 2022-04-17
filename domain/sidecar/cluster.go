package sidecar

import (
	adminapi "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	"istio.io/istio/pkg/util/protomarshal"
)

type Cluster struct {
	*adminapi.Clusters
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
