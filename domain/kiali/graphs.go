package kiali

import (
	"net/http"
)

const (
	NamespacesGraph = "/namespaces/graph"
)

/**
 *  appenders          : Comma-separated list of Appenders to run. Available appenders: [aggregateNode, deadNode, idleNode, istio, responseTime, securityPolicy, serviceEntry, sidecarsCheck], default run all appenders
 *  duration           : Query time-range duration (Golang string duration), default 10m
 *  graphType		   : Graph type. Available graph types: [app, service, versionedApp, workload], default workload
 *  boxBy              : Comma-separated list of desired node boxing. Available boxings: [app, cluster, namespace, none], default none
 *  includeIdleEdges   : Flag for including edges that have no request traffic for the time period, default false
 *  injectServiceNodes : Flag for injecting the requested service node between source and destination nodes, default false
 *  namespaces 		   : Comma-separated list of namespaces to include in the graph. The namespaces must be accessible to the client.
 *  queryTime          : Unix time (seconds) for query such that time range is [queryTime-duration..queryTime, default now
 **/
func (c *Client) GetNamespacesGraph(namespaces, appenders, duration, graphType, boxBy, includeIdleEdges,
	injectServiceNodes, queryTime string) map[string]interface{} {
	requestArgs := map[string]string{
		"appenders":          appenders,
		"duration":           duration,
		"graphType":          graphType,
		"boxBy":              boxBy,
		"includeIdleEdges":   includeIdleEdges,
		"injectServiceNodes": injectServiceNodes,
		"namespaces":         namespaces,
		"queryTime":          queryTime,
	}
	graphInfo := make(map[string]interface{})
	request, err := http.NewRequest(http.MethodGet, c.GetRequestUrl(NamespacesGraph, requestArgs), nil)
	if err != nil {
		domainLog.Errorf("build request err, err: %s", err)
		return graphInfo
	}
	err = c.DoRequest(request, &graphInfo)
	if err != nil {
		domainLog.Errorf("do request err, err: %s", err)
	}
	return graphInfo
}
