package istio

import (
	"context"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	informer "istio.io/client-go/pkg/listers/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type EnvoyFilter struct {
	*IstioClient
}

// Must be careful of envoyFilter !!!
func NewEnvoyFilter(cli *IstioClient) *EnvoyFilter {
	return &EnvoyFilter{cli}
}

func (e *EnvoyFilter) List(namespace string) []v1alpha3.EnvoyFilter {
	envoyFilterList := make([]v1alpha3.EnvoyFilter, 0)
	list, err := e.GetEnvoyFilterLister().EnvoyFilters(namespace).List(labels.Everything())
	if err != nil || len(list) == 0 {
		list, err := e.Clientset.NetworkingV1alpha3().EnvoyFilters(namespace).List(context.Background(), metav1.ListOptions{})
		if err == nil && list != nil {
			envoyFilterList = list.Items
		}
		return envoyFilterList
	}
	for _, envoyFilter := range list {
		envoyFilterList = append(envoyFilterList, *envoyFilter)
	}
	return envoyFilterList
}

func (e *EnvoyFilter) Get(namespace, envoyFilterName string) (*v1alpha3.EnvoyFilter, error) {
	envoyFilter, err := e.GetEnvoyFilterLister().EnvoyFilters(namespace).Get(envoyFilterName)
	if err != nil || envoyFilter == nil {
		return e.Clientset.NetworkingV1alpha3().EnvoyFilters(namespace).
			Get(context.Background(), envoyFilterName, metav1.GetOptions{})
	}
	return envoyFilter, err
}

func (e *EnvoyFilter) Create(envoyFilter *v1alpha3.EnvoyFilter) error {
	_, err := e.Clientset.NetworkingV1alpha3().EnvoyFilters(envoyFilter.Namespace).
		Create(context.Background(), envoyFilter, metav1.CreateOptions{})
	return err
}

func (e *EnvoyFilter) Delete(namespace, envoyFilterName string) error {
	return e.Clientset.NetworkingV1alpha3().EnvoyFilters(namespace).
		Delete(context.Background(), envoyFilterName, metav1.DeleteOptions{})
}

func (e *EnvoyFilter) Update(envoyFilter *v1alpha3.EnvoyFilter) error {
	_, err := e.Clientset.NetworkingV1alpha3().EnvoyFilters(envoyFilter.Namespace).
		Update(context.Background(), envoyFilter, metav1.UpdateOptions{})
	return err
}

func (e *EnvoyFilter) GetEnvoyFilterLister() informer.EnvoyFilterLister {
	return e.SharedInformerFactory.Networking().V1alpha3().EnvoyFilters().Lister()
}
