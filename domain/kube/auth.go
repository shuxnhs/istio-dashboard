package kube

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/shuxnhs/istio-dashboard/model"

	"k8s.io/client-go/rest"
)

func tryTlsAuth(config *rest.Config, kubeConfig *model.KubeConfig) error {
	var err error

	if strings.Trim(kubeConfig.K8sClusterAuthData, " ") != "" {
		config.CAData, err = base64.StdEncoding.DecodeString(kubeConfig.K8sClusterAuthData)
		if err != nil {
			config.CAData = nil
			return err
		}
	} else {
		//没有证书则不需要，跳过证书校验
		config.Insecure = true
	}

	//下面是认证的用户client-certificate-data和client-key-data信息
	if kubeConfig.K8sClientCertificateData != "" && kubeConfig.K8sClientKeyData != "" {
		config.CertData, err = base64.StdEncoding.DecodeString(kubeConfig.K8sClientCertificateData)
		if err != nil {
			config.CAData = nil
			config.CertData = nil
			return err
		}
		config.KeyData, err = base64.StdEncoding.DecodeString(kubeConfig.K8sClientKeyData)
		if err != nil {
			config.CAData = nil
			config.CertData = nil
			config.KeyData = nil
			return err
		}
		return nil
	}

	return fmt.Errorf("tls auth data error")
}

func tryBasicAuth(config *rest.Config, kubeConfig *model.KubeConfig) error {
	if kubeConfig.K8sAuthBasic == "" {
		return fmt.Errorf("basic auth data is empty")
	}
	config.Insecure = true
	usernameColonPassword, err := base64.StdEncoding.DecodeString(kubeConfig.K8sAuthBasic)
	if err != nil {
		return err
	}
	usernamePassword := strings.SplitN(string(usernameColonPassword), ":", 2)
	if len(usernamePassword) >= 2 {
		config.Username = usernamePassword[0]
		config.Password = usernamePassword[1]
	} else {
		return fmt.Errorf("basic auth data incorrect, decode username and password error")
	}
	return nil
}

func tryTokenAuth(config *rest.Config, kubeConfig *model.KubeConfig) error {
	if kubeConfig.K8sAuthToken == "" {
		return fmt.Errorf("token auth data is empty")
	}
	config.BearerToken = kubeConfig.K8sAuthToken
	config.Insecure = true
	return nil
}
