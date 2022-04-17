package kiali

import (
	"fmt"
	"net/http"

	"github.com/kiali/kiali/models"
)

const (
	Workloads        = "/namespaces/%s/workloads"
	WorkloadsHealth  = "/namespaces/%s/workloads/%s/health"
	WorkloadsMetrics = "/namespaces/%s/workloads/%s/metrics"
)

func (c *Client) GetWorkloads(namespace string) models.WorkloadList {
	workloads := models.WorkloadList{}
	request, err := http.NewRequest(http.MethodGet, c.GetRequestUrl(
		fmt.Sprintf(Workloads, namespace), map[string]string{}), nil)
	if err != nil {
		domainLog.Errorf("build request err, err: %s", err)
		return workloads
	}
	err = c.DoRequest(request, &workloads)
	if err != nil {
		domainLog.Errorf("do request err, err: %s", err)
	}
	return workloads
}

/**
 *  rateInterval : The rate interval used for fetching error rate, default 10m
 *  QueryTime    : The time to use for the prometheus query
 *  hType		 : The type of health, "app", "service" or "workload", default "workload"
 **/
func (c *Client) GetWorkloadsHealth(namespace, workload, rateInterval, QueryTime, hType string) models.NamespaceWorkloadHealth {
	health := make(models.NamespaceWorkloadHealth)
	request, err := http.NewRequest(http.MethodGet, c.GetRequestUrl(
		fmt.Sprintf(WorkloadsHealth, namespace, workload),
		map[string]string{
			"rateInterval": rateInterval,
			"QueryTime":    QueryTime,
			"type":         hType,
		}), nil)
	if err != nil {
		domainLog.Errorf("build request err, err: %s", err)
		return health
	}
	err = c.DoRequest(request, &health)
	if err != nil {
		domainLog.Errorf("do request err, err: %s", err)
	}
	return health
}

/**
 *  avg 			: Flag for fetching histogram average. Default is true.
 *  direction    	: Traffic direction: 'inbound' or 'outbound'. Default value : outbound
 *  duration		: Duration of the query period, in seconds. Default value : 1800
 *  rateFunc 		: Prometheus function used to calculate rate: 'rate' or 'irate'. Default value : rate
 *  rateInterval    : Interval used for rate and histogram calculation. Default value : 1m
 *  requestProtocol	: Desired request protocol for the telemetry: For example, 'http' or 'grpc'. Default value : all protocols
 *  reporter 		: Istio telemetry reporter: 'source' or 'destination'. Default value : source
 *  step    		: Step between [graph] datapoints, in seconds. Default value : 15
 *  version			: Filters metrics by the specified version.
 *  ArrayArgs
 *  byLabels 		: List of labels to use for grouping metrics (via Prometheus 'by' clause).
 *  filters    		: List of metrics to fetch. Fetch all metrics when empty. List entries are Kiali internal metric names. Default value : List []
 *  quantiles		: List of quantiles to fetch. Fetch no quantiles when empty. Ex: [0.5, 0.95, 0.99]. Default value : List []
 **/
func (c *Client) GetWorkloadsMetrics(namespace, workload, avg, direction, duration, rateFunc, rateInterval,
	requestProtocol, reporter, step, version string, byLabels, filters, quantiles []string) models.MetricsMap {
	metrics := make(models.MetricsMap)
	request, err := http.NewRequest(http.MethodGet,
		c.GetRequestUrl(fmt.Sprintf(WorkloadsMetrics, namespace, workload),
			map[string]string{
				"avg":             avg,
				"direction":       direction,
				"duration":        duration,
				"rateFunc":        rateFunc,
				"rateInterval":    rateInterval,
				"requestProtocol": requestProtocol,
				"reporter":        reporter,
				"step":            step,
				"version":         version,
			},
			map[string][]string{
				"byLabels":  byLabels,
				"filters":   filters,
				"quantiles": quantiles,
			}), nil)

	if err != nil {
		domainLog.Errorf("build request err, err: %s", err)
		return metrics
	}
	err = c.DoRequest(request, &metrics)
	if err != nil {
		domainLog.Errorf("do request err, err: %s", err)
	}
	return metrics
}
