package jaeger

import (
	"fmt"
	"net/http"
)

const (
	Services           = "/services"
	ServicesOperations = "/services/%s/operations"
	TracesList         = "/traces"
	TraceDetail        = "/traces/%s"
)

func (c *Client) GetServices() map[string]interface{} {
	var result map[string]interface{}
	request, err := http.NewRequest(http.MethodGet, c.GetRequestUrl(Services, map[string]string{}), nil)
	if err != nil {
		domainLog.Errorf("build request err, err: %s", err)
		return nil
	}
	err = c.DoRequest(request, &result)
	if err != nil {
		domainLog.Errorf("do request err, err: %s", err)
	}
	return result
}

func (c *Client) ServicesOperations(namespace, service string) map[string]interface{} {
	var result map[string]interface{}
	serviceParam := service
	if namespace != "" {
		serviceParam = service + "." + namespace
	}
	request, err := http.NewRequest(http.MethodGet, c.GetRequestUrl(
		fmt.Sprintf(ServicesOperations, serviceParam), map[string]string{}), nil)
	if err != nil {
		domainLog.Errorf("build request err, err: %s", err)
		return nil
	}
	err = c.DoRequest(request, &result)
	if err != nil {
		domainLog.Errorf("do request err, err: %s", err)
	}
	return result
}

func (c *Client) TraceList(service, namespace, operation, lookBack, tags, limit, start, end, maxDuration, minDuration string) map[string]interface{} {
	var result map[string]interface{}
	serviceParam := service
	if namespace != "" {
		serviceParam = service + "." + namespace
	}
	queryArgs := map[string]string{
		"service":     serviceParam,
		"operation":   operation,
		"lookback":    lookBack,
		"limit":       limit,
		"start":       start,
		"end":         end,
		"maxDuration": maxDuration,
		"minDuration": minDuration,
	}
	if tags != "" && tags != "{}" && tags != "{[]}" {
		queryArgs["tags"] = tags
	}
	request, err := http.NewRequest(http.MethodGet, c.GetRequestUrl(TracesList, queryArgs), nil)
	if err != nil {
		domainLog.Errorf("build request err, err: %s", err)
		return nil
	}
	err = c.DoRequest(request, &result)
	if err != nil {
		domainLog.Errorf("do request err, err: %s", err)
	}
	return result
}

func (c *Client) TraceDetail(traceId string) map[string]interface{} {
	var result map[string]interface{}
	request, err := http.NewRequest(http.MethodGet, c.GetRequestUrl(
		fmt.Sprintf(TraceDetail, traceId), map[string]string{}), nil)
	if err != nil {
		domainLog.Errorf("build request err, err: %s", err)
		return nil
	}
	err = c.DoRequest(request, &result)
	if err != nil {
		domainLog.Errorf("do request err, err: %s", err)
	}
	return result
}
