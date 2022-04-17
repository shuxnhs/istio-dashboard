package sidecar

import (
	adminapi "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	"github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
)

type EDSInfo map[string][]EDSClusterInfo

type EDSClusterInfo struct {
	Address            string `json:"address"`
	Port               int    `json:"port"`
	Status             string `json:"status"`
	FailedOutlierCheck bool   `json:"filedOutlierCheck"`
}

type EDS struct {
	Address            string `json:"address"`
	Port               int    `json:"port"`
	Cluster            string `json:"cluster"`
	Status             string `json:"status"`
	FailedOutlierCheck bool   `json:"filedOutlierCheck"`
}

func ClustersToEDSInfo(clusters *Cluster) EDSInfo {
	edsInfo := make(EDSInfo)
	for _, cluster := range clusters.ClusterStatuses {
		edsClusterInfos := make([]EDSClusterInfo, 0)
		for _, host := range cluster.HostStatuses {
			edsClusterInfos = append(edsClusterInfos, EDSClusterInfo{
				Address:            retrieveEndpointAddress(host),
				Port:               int(retrieveEndpointPort(host)),
				Status:             retrieveEndpointStatus(host).String(),
				FailedOutlierCheck: retrieveFailedOutlierCheck(host),
			})
		}
		edsInfo[cluster.Name] = edsClusterInfos
	}
	return edsInfo
}

func ClustersToEDS(clusters *Cluster) []EDS {
	eds := make([]EDS, 0)
	for _, cluster := range clusters.ClusterStatuses {
		for _, host := range cluster.HostStatuses {
			eds = append(eds, EDS{
				Address:            retrieveEndpointAddress(host),
				Port:               int(retrieveEndpointPort(host)),
				Cluster:            cluster.Name,
				Status:             retrieveEndpointStatus(host).String(),
				FailedOutlierCheck: retrieveFailedOutlierCheck(host),
			})
		}
	}

	return eds
}

func retrieveEndpointAddress(host *adminapi.HostStatus) string {
	addr := host.Address.GetSocketAddress()
	if addr != nil {
		return addr.Address
	}
	return "unix://" + host.Address.GetPipe().Path
}

func retrieveEndpointPort(l *adminapi.HostStatus) uint32 {
	addr := l.Address.GetSocketAddress()
	if addr != nil {
		return addr.GetPortValue()
	}
	return 0
}

func retrieveEndpointStatus(l *adminapi.HostStatus) corev3.HealthStatus {
	return l.HealthStatus.GetEdsHealthStatus()
}

func retrieveFailedOutlierCheck(l *adminapi.HostStatus) bool {
	return l.HealthStatus.GetFailedOutlierCheck()
}
