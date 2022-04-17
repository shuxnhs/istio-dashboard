package kube

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Deployment struct {
	cli *kubernetes.Clientset
}

func NewDeployment(cli *kubernetes.Clientset) *Deployment {
	return &Deployment{cli: cli}
}

func (d *Deployment) GetDeployment(namespace, deploymentName string) (*v1.Deployment, error) {
	return d.cli.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
}

func (d *Deployment) UpdateDeployment(namespace string, deployment *v1.Deployment) (*v1.Deployment, error) {
	return d.cli.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
}
