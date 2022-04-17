package kiali

import (
	"fmt"
	"net/http"

	"github.com/kiali/kiali/models"
)

const (
	Namespaces        = "/namespaces"
	NamespacesHealth  = "/namespaces/%s/health"
	NamespacesMetrics = "/namespaces/%s/metrics"
)

func (c *Client) GetNamespaces() []models.Namespace {
	namespaces := make([]models.Namespace, 0)
	request, err := http.NewRequest(http.MethodGet, c.GetRequestUrl(Namespaces, map[string]string{}), nil)
	if err != nil {
		domainLog.Errorf("build request err, err: %s", err)
		return namespaces
	}
	err = c.DoRequest(request, &namespaces)
	if err != nil {
		domainLog.Errorf("do request err, err: %s", err)
	}
	return namespaces
}

/**
 *  rateInterval : The rate interval used for fetching error rate, default 10m
 *  QueryTime    : The time to use for the prometheus query
 *  hType		 : The type of health, "app", "service" or "workload", default "app"
 **/
func (c *Client) GetNamespaceHealth(namespace, rateInterval, queryTime, hType string) models.NamespaceAppHealth {
	health := make(models.NamespaceAppHealth)
	request, err := http.NewRequest(
		http.MethodGet,
		c.GetRequestUrl(
			fmt.Sprintf(NamespacesHealth, namespace),
			map[string]string{
				"rateInterval": rateInterval,
				"QueryTime":    queryTime,
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

func (c *Client) GetNamespaceMetrics(namespace string) models.MetricsMap {
	metrics := make(models.MetricsMap)
	request, err := http.NewRequest(
		http.MethodGet,
		c.GetRequestUrl(fmt.Sprintf(NamespacesMetrics, namespace), map[string]string{}),
		nil)

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
