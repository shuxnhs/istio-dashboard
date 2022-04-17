package kube

import (
	"context"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Namespace struct {
	cli *kubernetes.Clientset
}

func NewNamespace(cli *kubernetes.Clientset) *Namespace {
	return &Namespace{cli: cli}
}

func (n *Namespace) GetNamespaceObject(namespace string) (*v1.Namespace, error) {
	return n.cli.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
}

func (n *Namespace) UpdateNamespaceObject(namespace *v1.Namespace) (*v1.Namespace, error) {
	return n.cli.CoreV1().Namespaces().Update(context.Background(), namespace, metav1.UpdateOptions{})
}

func (n *Namespace) GetNamespaceLabel(namespace string) (map[string]string, error) {
	namespaceInfo, err := n.GetNamespaceObject(namespace)
	if err != nil {
		return map[string]string{}, err
	}
	return namespaceInfo.Labels, nil
}

func (n *Namespace) ListNamespaceByLabel(label string) (*v1.NamespaceList, error) {
	return n.cli.CoreV1().Namespaces().
		List(context.Background(), metav1.ListOptions{LabelSelector: label})
}
