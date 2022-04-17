package istio

import (
	"context"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	informer "istio.io/client-go/pkg/listers/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type VirtualService struct {
	*IstioClient
}

func NewVirtualService(cli *IstioClient) *VirtualService {
	return &VirtualService{cli}
}

func (v *VirtualService) Get(namespace, virtualServiceName string) (*v1alpha3.VirtualService, error) {
	virtualService, err := v.GetVirtualServiceLister().VirtualServices(namespace).Get(virtualServiceName)
	if virtualService == nil {
		return v.Clientset.NetworkingV1alpha3().VirtualServices(namespace).
			Get(context.Background(), virtualServiceName, metav1.GetOptions{})
	}
	return virtualService, err
}

func (v *VirtualService) List(namespace string, opts metav1.ListOptions) []v1alpha3.VirtualService {
	virtualServiceList := make([]v1alpha3.VirtualService, 0)
	list, err := v.GetVirtualServiceLister().VirtualServices(namespace).List(labels.Everything())
	if err != nil || len(list) == 0 {
		list, err := v.Clientset.NetworkingV1alpha3().VirtualServices(namespace).List(context.Background(), opts)
		if err == nil && list != nil {
			virtualServiceList = list.Items
		}
		return virtualServiceList
	}
	for _, virtualService := range list {
		virtualServiceList = append(virtualServiceList, *virtualService)
	}
	return virtualServiceList
}

func (v *VirtualService) Create(vs *v1alpha3.VirtualService) error {
	_, k8sErr := v.Clientset.NetworkingV1alpha3().VirtualServices(vs.Namespace).Create(context.Background(), vs, metav1.CreateOptions{})
	if k8sErr != nil {
		return k8sErr
	}
	return nil
}

func (v *VirtualService) Update(vs *v1alpha3.VirtualService) error {
	_, k8sErr := v.Clientset.NetworkingV1alpha3().VirtualServices(vs.Namespace).Update(context.Background(), vs, metav1.UpdateOptions{})
	if k8sErr != nil {
		return k8sErr
	}
	return nil
}

func (v *VirtualService) Delete(namespace, name string) error {
	k8sErr := v.Clientset.NetworkingV1alpha3().VirtualServices(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if k8sErr != nil {
		if errors.IsNotFound(k8sErr) {
			return nil
		}
		return k8sErr
	}

	return nil
}

func (v *VirtualService) DoCreateOrUpdate(virtualService *v1alpha3.VirtualService) error {
	exist := true
	oldVirtualService, err := v.Get(virtualService.Namespace, virtualService.Name)
	if err != nil {
		if errors.IsNotFound(err) {
			exist = false
		} else {
			return err
		}
	}

	if !exist {
		return v.Create(virtualService)
	} else {
		oldVirtualService.Spec = virtualService.Spec
		return v.Update(oldVirtualService)
	}
}

func (v *VirtualService) GetVirtualServiceLister() informer.VirtualServiceLister {
	return v.SharedInformerFactory.Networking().V1alpha3().VirtualServices().Lister()
}
