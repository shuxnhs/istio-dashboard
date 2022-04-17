package kube

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type StatefulSet struct {
	cli *kubernetes.Clientset
}

func NewStatefulSet(cli *kubernetes.Clientset) *StatefulSet {
	return &StatefulSet{cli: cli}
}

func (s *StatefulSet) GetStatefulSet(namespace, statefulSetName string) (*v1.StatefulSet, error) {
	return s.cli.AppsV1().StatefulSets(namespace).Get(context.Background(), statefulSetName, metav1.GetOptions{})
}

func (s *StatefulSet) UpdateStatefulSet(namespace string, statefulSet *v1.StatefulSet) (*v1.StatefulSet, error) {
	return s.cli.AppsV1().StatefulSets(namespace).Update(context.Background(), statefulSet, metav1.UpdateOptions{})
}
