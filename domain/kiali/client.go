package kiali

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/shuxnhs/istio-dashboard/domain/kube"
	"github.com/shuxnhs/istio-dashboard/model"

	"istio.io/pkg/log"
	"k8s.io/client-go/rest"
)

var domainLog = log.RegisterScope("kiali-domain", "kiali-domain debugging", 0)

const KialiPath = "/api/v1/namespaces/istio-system/services/kiali:http/proxy/kiali/api"

type Client struct {
	kialiPath  string
	kubeConfig *model.KubeConfig
	*rest.RESTClient
}

func NewKialiClient(kubeConfig *model.KubeConfig) *Client {
	restClient := kube.NewKubernetesRestClient(kubeConfig)
	if restClient != nil {
		return &Client{
			kialiPath:  kubeConfig.KialiPath,
			kubeConfig: kubeConfig,
			RESTClient: restClient,
		}
	}
	return nil
}

func (c *Client) GetRequestUrl(apiName string, queryArgs map[string]string, queryArrayArgs ...map[string][]string) string {
	uri := c.kubeConfig.K8sHost + c.kialiPath + apiName + c.buildQueryString(queryArgs)
	if len(queryArrayArgs) > 0 {
		uri += "&"
		for k, v := range queryArrayArgs[0] {
			for _, n := range v {
				uri += k + "[]=" + url.QueryEscape(n) + "&"
			}
		}
		uri = strings.TrimSuffix(uri, "&")
	}
	domainLog.Infof("kiali client request url: %s", uri)
	return uri
}

func (c *Client) buildQueryString(queryArgs map[string]string) string {
	if len(queryArgs) > 0 {
		query := "?"
		for k, v := range queryArgs {
			query = query + k + "=" + url.QueryEscape(v) + "&"
		}
		return strings.TrimSuffix(query, "&")
	}
	return ""
}

func (c *Client) DoRequest(request *http.Request, data interface{}) error {
	response, err := c.RESTClient.Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return err
	}
	return nil
}
