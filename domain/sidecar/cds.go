package sidecar

import (
	"fmt"
	"strings"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/networking/util"
	"istio.io/istio/pkg/config/host"
)

type TrafficDirection string

const (
	TrafficDirectionInbound  TrafficDirection = "inbound"
	TrafficDirectionOutbound TrafficDirection = "outbound"
)

type CDS struct {
	FQDN            host.Name              `json:"fqdn"`
	Port            int                    `json:"port"`
	Subset          string                 `json:"subset"`
	Direction       model.TrafficDirection `json:"direction"`
	Type            string                 `json:"type"`
	DestinationRule string                 `json:"destinationRule"`
}

func ClustersToCDS(configDump *ConfigDump) []CDS {
	cds := make([]CDS, 0)
	clusters, err := configDump.GetCluster()
	if err != nil {
		return cds
	}
	for _, cluster := range clusters {
		if len(strings.Split(cluster.GetName(), "|")) > 3 {
			direction, subset, fqdn, port := model.ParseSubsetKey(cluster.GetName())
			if subset == "" {
				subset = "-"
			}
			cds = append(cds, CDS{
				FQDN:            fqdn,
				Port:            port,
				Subset:          subset,
				Direction:       direction,
				Type:            cluster.GetType().String(),
				DestinationRule: mdToDestinationRule(cluster.GetMetadata()),
			})
		} else {
			cds = append(cds, CDS{
				FQDN:            host.Name(cluster.GetName()),
				Port:            0,
				Subset:          "-",
				Direction:       "-",
				Type:            cluster.GetType().String(),
				DestinationRule: mdToDestinationRule(cluster.GetMetadata()),
			})
		}
	}
	return cds
}

func mdToDestinationRule(metadata *core.Metadata) string {
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
	if strings.HasPrefix(config.GetStringValue(), "/apis/networking.istio.io/v1alpha3/namespaces/") {
		pieces := strings.Split(config.GetStringValue(), "/")
		if len(pieces) != 8 {
			return ""
		}
		return fmt.Sprintf("%s.%s", pieces[7], pieces[5])
	}
	return "<unknown>"
}
