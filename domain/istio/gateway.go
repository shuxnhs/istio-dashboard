package istio

import (
	"context"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	informer "istio.io/client-go/pkg/listers/networking/v1alpha3"
	"istio.io/pkg/log"
	kerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type Gateway struct {
	*IstioClient
}

func NewGateway(cli *IstioClient) *Gateway {
	return &Gateway{cli}
}

func (g *Gateway) Get(namespace, gatewayName string) (*v1alpha3.Gateway, error) {
	gateway, err := g.GetGatewayLister().Gateways(namespace).Get(gatewayName)
	if gateway == nil {
		return g.Clientset.NetworkingV1alpha3().Gateways(namespace).
			Get(context.Background(), gatewayName, metav1.GetOptions{})
	}
	return gateway, err
}

func (g *Gateway) List(namespace string, label map[string]string) ([]*v1alpha3.Gateway, error) {
	selector := labels.Everything()
	if len(label) != 0 {
		selector = labels.SelectorFromSet(label)
	}

	gateways, err := g.GetGatewayLister().Gateways(namespace).List(selector)
	if gateways == nil {
		gws, err := g.Clientset.NetworkingV1alpha3().Gateways(namespace).
			List(context.Background(), metav1.ListOptions{LabelSelector: selector.String()})
		if err != nil {
			return nil, err
		}

		for i := range gws.Items {
			gateways = append(gateways, &gws.Items[i])
		}
	}

	return gateways, err
}

//informer查询列表
func (g *Gateway) ListCache(namespace string, label map[string]string) ([]*v1alpha3.Gateway, error) {
	selector := labels.Everything()
	if len(label) != 0 {
		selector = labels.SelectorFromSet(label)
	}

	return g.GetGatewayLister().Gateways(namespace).List(selector)
}

func (g *Gateway) Create(gateway *v1alpha3.Gateway) error {
	_, k8sErr := g.Clientset.NetworkingV1alpha3().Gateways(gateway.Namespace).Create(context.Background(), gateway, metav1.CreateOptions{})
	if k8sErr != nil {
		return k8sErr
	}
	return nil
}

func (g *Gateway) Update(gateway *v1alpha3.Gateway) error {
	_, k8sErr := g.Clientset.NetworkingV1alpha3().Gateways(gateway.Namespace).Update(context.Background(), gateway, metav1.UpdateOptions{})
	if k8sErr != nil {
		return k8sErr
	}
	return nil
}

func (g *Gateway) Delete(namespace, gatewayName string) error {
	k8sErr := g.Clientset.NetworkingV1alpha3().Gateways(namespace).Delete(context.Background(), gatewayName, metav1.DeleteOptions{})
	if k8sErr != nil {
		if kerror.IsNotFound(k8sErr) {
			return nil
		}
		log.Error(k8sErr, namespace, gatewayName)
		return k8sErr
	}
	return nil
}

func (g *Gateway) GetGatewayLister() informer.GatewayLister {
	return g.SharedInformerFactory.Networking().V1alpha3().Gateways().Lister()
}
