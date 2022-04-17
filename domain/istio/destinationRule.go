package istio

import (
	"context"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	informer "istio.io/client-go/pkg/listers/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type DestinationRule struct {
	*IstioClient
}

func NewDestinationRule(cli *IstioClient) *DestinationRule {
	return &DestinationRule{cli}
}

func (d *DestinationRule) List(namespace string) []v1alpha3.DestinationRule {
	destinationRuleList := make([]v1alpha3.DestinationRule, 0)
	list, err := d.GetDestinationRuleLister().DestinationRules(namespace).List(labels.Everything())
	if err != nil || len(list) == 0 {
		list, err := d.Clientset.NetworkingV1alpha3().DestinationRules(namespace).List(context.Background(), metav1.ListOptions{})
		if err == nil && list != nil {
			destinationRuleList = list.Items
		}
		return destinationRuleList
	}
	for _, destinationRule := range list {
		destinationRuleList = append(destinationRuleList, *destinationRule)
	}
	return destinationRuleList
}

func (d *DestinationRule) Get(namespace, destinationRuleName string) (*v1alpha3.DestinationRule, error) {
	destinationRule, err := d.GetDestinationRuleLister().DestinationRules(namespace).Get(destinationRuleName)
	if destinationRule == nil {
		return d.Clientset.NetworkingV1alpha3().DestinationRules(namespace).
			Get(context.Background(), destinationRuleName, metav1.GetOptions{})
	}
	return destinationRule, err
}

func (d *DestinationRule) Create(destinationRule *v1alpha3.DestinationRule) error {
	_, err := d.Clientset.NetworkingV1alpha3().DestinationRules(destinationRule.Namespace).
		Create(context.Background(), destinationRule, metav1.CreateOptions{})
	return err
}

func (d *DestinationRule) Delete(namespace, destinationRuleName string) error {
	return d.Clientset.NetworkingV1alpha3().DestinationRules(namespace).
		Delete(context.Background(), destinationRuleName, metav1.DeleteOptions{})
}

func (d *DestinationRule) Update(destinationRule *v1alpha3.DestinationRule) error {
	_, err := d.Clientset.NetworkingV1alpha3().DestinationRules(destinationRule.Namespace).
		Update(context.Background(), destinationRule, metav1.UpdateOptions{})
	return err
}

func (d *DestinationRule) DoCreateOrUpdate(destinationRule *v1alpha3.DestinationRule) error {
	oldDstinationRule, err := d.Get(destinationRule.Namespace, destinationRule.Name)
	exist := true
	if err != nil {
		if errors.IsNotFound(err) {
			exist = false
		} else {
			return err
		}
	}
	if !exist {
		// create
		return d.Create(destinationRule)
	} else {
		// update
		oldDstinationRule.Spec = destinationRule.Spec
		return d.Update(oldDstinationRule)
	}
}

func (d *DestinationRule) GetDestinationRuleLister() informer.DestinationRuleLister {
	return d.SharedInformerFactory.Networking().V1alpha3().DestinationRules().Lister()
}
