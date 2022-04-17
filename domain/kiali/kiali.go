package kiali

import (
	"net/http"

	"github.com/kiali/kiali/status"
)

const (
	Config      = "/config"
	Status      = "/status"
	IstioStatus = "/istio/status"
)

//func (c *Client) GetKialiConfig() handlers.PublicConfig {
//	config := handlers.PublicConfig{}
//	request, err := http.NewRequest(http.MethodGet, c.GetRequestUrl(Config, map[string]string{}), nil)
//	if err != nil {
//		domainLog.Errorf("build request err, err: %s", err)
//		return config
//	}
//	err = c.DoRequest(request, &config)
//	if err != nil {
//		domainLog.Errorf("do request err, err: %s", err)
//	}
//	return config
//}

func (c *Client) GetKialiStatus() status.StatusInfo {
	kialiStatus := status.StatusInfo{}
	request, err := http.NewRequest(http.MethodGet, c.GetRequestUrl(Status, map[string]string{}), nil)
	if err != nil {
		domainLog.Errorf("build request err, err: %s", err)
		return kialiStatus
	}
	err = c.DoRequest(request, &kialiStatus)
	if err != nil {
		domainLog.Errorf("do request err, err: %s", err)
	}
	return kialiStatus
}

//func (c *Client) GetKialiIstioStatus() business.IstioComponentStatus {
//	istioStatus := make(business.IstioComponentStatus, 0)
//	request, err := http.NewRequest(http.MethodGet, c.GetRequestUrl(IstioStatus, map[string]string{}), nil)
//	if err != nil {
//		domainLog.Errorf("build request err, err: %s", err)
//		return istioStatus
//	}
//	err = c.DoRequest(request, &istioStatus)
//	if err != nil {
//		domainLog.Errorf("do request err, err: %s", err)
//	}
//	return istioStatus
//}
