package sidecar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/shuxnhs/istio-dashboard/domain/kube"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const istioNamespace = "istio-system"

type Sidecar struct {
	config  *rest.Config
	cli     *kubernetes.Clientset
	restCli *rest.RESTClient
}

func NewSidecar(config *rest.Config) *Sidecar {
	return &Sidecar{config: config, cli: kube.NewClientSet(config), restCli: kube.NewRestClient(config)}
}

func (s *Sidecar) Check(namespace, pod string) {
	//path := "config_dump"
	//config, err := s.EnvoyDo(context.TODO(), pod, namespace, "GET", path)
	//if err != nil {
	//	return nil, err
	//}
	//
	//path = fmt.Sprintf("/debug/config_dump?proxyID=%s.%s", pod, namespace)
	//istiodDumps, err := s.AllDiscoveryDo(context.TODO(), istioNamespace, path)

	//statuses, err := s.AllDiscoveryDo(context.TODO(), istioNamespace, "/debug/syncz")
	//if err != nil {
	//
	//}

}

func (s *Sidecar) GetEDS(namespace, pod string) (EDSInfo, error) {
	path := "clusters?format=json"
	config, err := s.EnvoyDo(context.TODO(), pod, namespace, "GET", path)
	if err != nil {
		return nil, err
	}
	cluster, err := s.ClustersFormat(config)
	if err != nil {
		return nil, err
	}
	return ClustersToEDSInfo(cluster), nil
}

func (s *Sidecar) ClustersFormat(config []byte) (*Cluster, error) {
	cluster := &Cluster{}
	err := json.Unmarshal(config, cluster)
	if err != nil {
		return nil, err
	}
	return cluster, nil

}

func (s *Sidecar) EnvoyDo(ctx context.Context, podName, podNamespace, method, path string) ([]byte, error) {
	return s.portForwardRequest(ctx, podName, podNamespace, method, path, 15000)
}

func (s *Sidecar) AllDiscoveryDo(ctx context.Context, istiodNamespace, path string) (map[string][]byte, error) {
	istiods, err := s.cli.CoreV1().Pods(istiodNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: fields.SelectorFromSet(map[string]string{
			"labelSelector": "app=istiod",
			"fieldSelector": "status.phase=Running",
		}).String(),
	})
	if err != nil {
		return nil, err
	}
	if len(istiods.Items) == 0 {
		return nil, errors.New("unable to find any Istiod instances")
	}

	result := map[string][]byte{}
	for _, istiod := range istiods.Items {
		res, err := s.portForwardRequest(ctx, istiod.Name, istiod.Namespace, http.MethodGet, path, 15014)
		if err != nil {
			return nil, err
		}
		if len(res) > 0 {
			result[istiod.Name] = res
		}
	}
	// If any Discovery servers responded, treat as a success
	if len(result) > 0 {
		return result, nil
	}
	return nil, nil
}

func (s *Sidecar) portForwardRequest(ctx context.Context, podName, podNamespace, method, path string, port int) ([]byte, error) {
	formatError := func(err error) error {
		return fmt.Errorf("failure running port forward process: %v", err)
	}

	fw, err := kube.NewPortForwarder(s.config, podName, podNamespace, "127.0.0.1", 0, port)
	if err != nil {
		return nil, err
	}
	if err = fw.Start(); err != nil {
		return nil, formatError(err)
	}
	defer fw.Close()
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s/%s", fw.Address(), path), nil)
	if err != nil {
		return nil, formatError(err)
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, formatError(err)
	}
	defer resp.Body.Close()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, formatError(err)
	}

	return out, nil
}
