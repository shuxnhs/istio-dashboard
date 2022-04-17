package kube

import (
	"github.com/shuxnhs/istio-dashboard/model"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetConfigStoreKubeConfig(kubeConfig *model.KubeConfig) *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags(kubeConfig.K8sHost, "")
	if err != nil {
		domainLog.Errorf("build config err: %s", err)
		return nil
	}

	switch kubeConfig.K8sAuthType {
	case model.K8sAuthTypeTLS:
		if err := tryTlsAuth(config, kubeConfig); err != nil {
			domainLog.Errorf("build config and tls auth err: %s", err)
			return nil
		}
	case model.K8sAuthTypeBASIC:
		if err := tryBasicAuth(config, kubeConfig); err != nil {
			domainLog.Errorf("build config and b basic auth err: %s", err)
			return nil
		}
	case model.K8sAuthTypeTOKEN:
		if err := tryTokenAuth(config, kubeConfig); err != nil {
			domainLog.Errorf("build config and b basic token err: %s", err)
			return nil
		}
	case model.K8sAuthInCLUSTER:
		inClusterCfg, err := rest.InClusterConfig()
		if err == nil {
			return inClusterCfg
		}
	}

	// 支持不认证
	config.Insecure = true
	return config
}

func NewKubernetesClientSet(kubeConfig *model.KubeConfig) *kubernetes.Clientset {
	return NewClientSet(GetConfigStoreKubeConfig(kubeConfig))
}

func NewKubernetesRestClient(kubeConfig *model.KubeConfig) *rest.RESTClient {
	return NewRestClient(GetConfigStoreKubeConfig(kubeConfig))
}
