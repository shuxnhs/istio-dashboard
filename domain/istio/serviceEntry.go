package istio

import (
	"context"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	informer "istio.io/client-go/pkg/listers/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type ServiceEntry struct {
	*IstioClient
}

func NewServiceEntry(cli *IstioClient) *ServiceEntry {
	return &ServiceEntry{cli}
}

func (s *ServiceEntry) List(namespace string) []v1alpha3.ServiceEntry {
	serviceEntryList := make([]v1alpha3.ServiceEntry, 0)
	list, err := s.GetServiceEntryLister().ServiceEntries(namespace).List(labels.Everything())
	if err != nil || len(list) == 0 {
		list, err := s.Clientset.NetworkingV1alpha3().ServiceEntries(namespace).List(context.Background(), metav1.ListOptions{})
		if err == nil && list != nil {
			serviceEntryList = list.Items
		}
		return serviceEntryList
	}
	for _, serviceEntry := range list {
		serviceEntryList = append(serviceEntryList, *serviceEntry)
	}
	return serviceEntryList
}

func (s *ServiceEntry) Get(namespace, serviceEntryName string) (*v1alpha3.ServiceEntry, error) {
	serviceEntry, err := s.GetServiceEntryLister().ServiceEntries(namespace).Get(serviceEntryName)
	if serviceEntry == nil {
		return s.Clientset.NetworkingV1alpha3().ServiceEntries(namespace).
			Get(context.Background(), serviceEntryName, metav1.GetOptions{})
	}
	return serviceEntry, err
}

func (s *ServiceEntry) Create(serviceEntry *v1alpha3.ServiceEntry) error {
	_, err := s.Clientset.NetworkingV1alpha3().ServiceEntries(serviceEntry.Namespace).
		Create(context.Background(), serviceEntry, metav1.CreateOptions{})
	return err
}

func (s *ServiceEntry) Delete(namespace, serviceEntryName string) error {
	return s.Clientset.NetworkingV1alpha3().ServiceEntries(namespace).
		Delete(context.Background(), serviceEntryName, metav1.DeleteOptions{})
}

func (s *ServiceEntry) Update(serviceEntry *v1alpha3.ServiceEntry) error {
	_, err := s.Clientset.NetworkingV1alpha3().ServiceEntries(serviceEntry.Namespace).
		Update(context.Background(), serviceEntry, metav1.UpdateOptions{})
	return err
}

func (s *ServiceEntry) GetServiceEntryLister() informer.ServiceEntryLister {
	return s.SharedInformerFactory.Networking().V1alpha3().ServiceEntries().Lister()
}
