package kube

import (
	"github.com/shuxnhs/istio-dashboard/model"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	NamespaceInjectionLabel    = "istio-injection"
	NamespaceInjectionEnable   = "enabled"
	NamespaceInjectionDisabled = "disabled"
	PodInjectionAnnotations    = "sidecar.istio.io/inject"
)

type InjectionManager struct {
	*kubernetes.Clientset
	*Namespace
	*Deployment
	*StatefulSet
}

func NewInjectionManager(kubeConfig *model.KubeConfig) *InjectionManager {
	cli := NewKubernetesClientSet(kubeConfig)
	return &InjectionManager{
		Clientset:   cli,
		Namespace:   NewNamespace(cli),
		Deployment:  NewDeployment(cli),
		StatefulSet: NewStatefulSet(cli),
	}
}

// 老版本需要先给命名空间打上istio-injection=enabled的label，然后根据pod的annoations来控制注入

// 检查命名空间是否开启自动注入
func (i *InjectionManager) CheckNamespaceAutoInjection(namespace string) bool {
	namespcaeLabels, err := i.Namespace.GetNamespaceLabel(namespace)
	if err != nil {
		return false
	}
	for k, v := range namespcaeLabels {
		if k == NamespaceInjectionLabel && v == NamespaceInjectionEnable {
			return true
		}
	}
	return false
}

// 为命名空间开启自动注入
func (i *InjectionManager) InjectNamespace(namespace string) bool {
	return i.setNamespaceInject(namespace, true)
}

func (i *InjectionManager) UnInjectNamespace(namespace string) bool {
	return i.setNamespaceInject(namespace, false)
}

// 获取所有开启了自动注入的命名空间
func (i *InjectionManager) ListInjectNamespace() *v1.NamespaceList {
	namespaces, err := i.Namespace.ListNamespaceByLabel(NamespaceInjectionLabel + "=" + NamespaceInjectionEnable)
	if err != nil {
		return &v1.NamespaceList{}
	}
	return namespaces
}

// pod自动注入, 只支持deployment的资源注入
func (i *InjectionManager) OldInjectInDeployment(namespace, deploymentName string) bool {
	deployment, err := i.Deployment.GetDeployment(namespace, deploymentName)
	if err != nil || deployment == nil {
		return false
	}
	return i.deploymentSidecarInject(namespace, deployment, true)
}

// pod取消注入
func (i *InjectionManager) OldInjectOutDeployment(namespace, deploymentName string) bool {
	deployment, err := i.Deployment.GetDeployment(namespace, deploymentName)
	if err != nil || deployment == nil {
		return false
	}
	return i.deploymentSidecarInject(namespace, deployment, false)
}

func (i *InjectionManager) setNamespaceInject(namespace string, enable bool) bool {
	namespaceObject, err := i.Namespace.GetNamespaceObject(namespace)
	if err != nil {
		return false
	}

	if len(namespaceObject.Labels) == 0 {
		namespaceObject.SetLabels(map[string]string{})
	}

	if enable {
		namespaceObject.Labels[NamespaceInjectionLabel] = NamespaceInjectionEnable
	} else {
		namespaceObject.Labels[NamespaceInjectionLabel] = NamespaceInjectionDisabled
	}

	if _, err = i.Namespace.UpdateNamespaceObject(namespaceObject); err != nil {
		return false
	}
	return true
}

// only support inject deployment
func (i *InjectionManager) deploymentSidecarInject(namespace string, deployment *appsv1.Deployment, inject bool) bool {
	// https://preliminary.istio.io/latest/zh/docs/ops/common-problems/injection/
	if deployment.Spec.Template.Spec.HostNetwork {
		return false
	}
	if len(deployment.Spec.Template.Annotations) == 0 {
		deployment.Spec.Template.Annotations = make(map[string]string)
	}
	if inject {
		deployment.Spec.Template.Annotations[PodInjectionAnnotations] = "true"
	} else {
		deployment.Spec.Template.Annotations[PodInjectionAnnotations] = "false"
	}
	if _, err := i.Deployment.UpdateDeployment(namespace, deployment); err == nil {
		return true
	}
	return false
}

// istio1.9+版本可以通过对pod的label进行控制，不用先给namespace创建对应的label
// github：https://github.com/istio/istio/issues/32388
func (i *InjectionManager) NewInjectInDeployment(namespace, deploymentName string) bool {
	deployment, err := i.Deployment.GetDeployment(namespace, deploymentName)
	if err != nil || deployment == nil {
		return false
	}
	return i.newSidecarInject(namespace, deployment, nil, true)
}

func (i *InjectionManager) NewInjectOutDeployment(namespace, deploymentName string) bool {
	deployment, err := i.Deployment.GetDeployment(namespace, deploymentName)
	if err != nil || deployment == nil {
		return false
	}
	return i.newSidecarInject(namespace, deployment, nil, false)
}

func (i *InjectionManager) NewInjectStatefulSet(namespace, statefulSetName string) bool {
	statefulSet, err := i.StatefulSet.GetStatefulSet(namespace, statefulSetName)
	if err != nil || statefulSet == nil {
		return false
	}
	return i.newSidecarInject(namespace, nil, statefulSet, true)
}

func (i *InjectionManager) NewInjectOutStatefulSet(namespace, statefulSetName string) bool {
	statefulSet, err := i.StatefulSet.GetStatefulSet(namespace, statefulSetName)
	if err != nil || statefulSet == nil {
		return false
	}
	return i.newSidecarInject(namespace, nil, statefulSet, false)
}

func (i *InjectionManager) newSidecarInject(namespace string, deployment *appsv1.Deployment, statefulSet *appsv1.StatefulSet, inject bool) bool {
	// https://preliminary.istio.io/latest/zh/docs/ops/common-problems/injection/
	if deployment != nil {
		if deployment.Spec.Template.Spec.HostNetwork {
			return false
		}
		if len(deployment.Spec.Template.Labels) == 0 {
			deployment.Spec.Template.Labels = make(map[string]string)
		}
		if inject {
			deployment.Spec.Template.Labels[PodInjectionAnnotations] = "true"
		} else {
			deployment.Spec.Template.Labels[PodInjectionAnnotations] = "false"
		}
		if _, err := i.Deployment.UpdateDeployment(namespace, deployment); err == nil {
			return true
		}
	}

	if statefulSet != nil {
		if statefulSet.Spec.Template.Spec.HostNetwork {
			return false
		}
		if len(statefulSet.Spec.Template.Labels) == 0 {
			statefulSet.Spec.Template.Labels = make(map[string]string)
		}
		if inject {
			statefulSet.Spec.Template.Labels[PodInjectionAnnotations] = "true"
		} else {
			statefulSet.Spec.Template.Labels[PodInjectionAnnotations] = "false"
		}
		if _, err := i.StatefulSet.UpdateStatefulSet(namespace, statefulSet); err == nil {
			return true
		}
	}

	return false
}
