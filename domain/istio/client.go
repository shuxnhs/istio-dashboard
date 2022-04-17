package istio

import (
	"time"

	"github.com/shuxnhs/istio-dashboard/domain/kube"
	"github.com/shuxnhs/istio-dashboard/model"

	"istio.io/client-go/pkg/clientset/versioned"
	"istio.io/client-go/pkg/informers/externalversions"
	"istio.io/pkg/log"
	"k8s.io/client-go/kubernetes"
)

const (
	IstioNamespace           = "istio-system"
	defaultIstioResyncPeriod = 300 * time.Second
)

var domainLog = log.RegisterScope("istio-domain", "istio-domain debugging", 0)

type IstioClient struct {
	stopChan chan struct{}
	kubeCli  *kubernetes.Clientset
	*versioned.Clientset
	externalversions.SharedInformerFactory
}

func NewIstioClientSet(kubeConfig *model.KubeConfig) *IstioClient {
	config := kube.GetConfigStoreKubeConfig(kubeConfig)
	istioClientSet, err := versioned.NewForConfig(config)
	if err != nil {
		domainLog.Errorf("new client client err: %s, config: %#v", err, config)
		return nil
	}
	istioClient := &IstioClient{
		stopChan:              make(chan struct{}),
		kubeCli:               kube.NewKubernetesClientSet(kubeConfig),
		Clientset:             istioClientSet,
		SharedInformerFactory: externalversions.NewSharedInformerFactory(istioClientSet, defaultIstioResyncPeriod),
	}
	for _, istioResource := range KindToIstioResourceSlice {
		genericInformer, err := istioClient.SharedInformerFactory.ForResource(istioResource)
		if err != nil {
			domainLog.Errorf("new sharedInformerFactory for resource %#v, err: %s", istioResource, err)
			return istioClient
		}
		go genericInformer.Informer().Run(istioClient.stopChan)
	}
	istioClient.SharedInformerFactory.Start(istioClient.stopChan)
	return istioClient
}

func (i *IstioClient) GetIstioVersion() {

}

func (i *IstioClient) CheckIstio() {

}
