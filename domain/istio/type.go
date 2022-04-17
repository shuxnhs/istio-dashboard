package istio

import (
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	ResourceNameGateway         = "gateways"
	ResourceNameDestinationRule = "destinationrules"
	ResourceNameVirtualService  = "virtualservices"
	ResourceNameServiceEntry    = "serviceentries"
	ResourceNameEnvoyFilter     = "envoyfilters"
)

var KindToIstioResourceSlice = []schema.GroupVersionResource{
	schema.GroupVersionResource{
		Group:    v1alpha3.GroupName,
		Version:  v1alpha3.SchemeGroupVersion.Version,
		Resource: ResourceNameGateway,
	},
	schema.GroupVersionResource{
		Group:    v1alpha3.GroupName,
		Version:  v1alpha3.SchemeGroupVersion.Version,
		Resource: ResourceNameDestinationRule,
	},
	schema.GroupVersionResource{
		Group:    v1alpha3.GroupName,
		Version:  v1alpha3.SchemeGroupVersion.Version,
		Resource: ResourceNameVirtualService,
	},
	schema.GroupVersionResource{
		Group:    v1alpha3.GroupName,
		Version:  v1alpha3.SchemeGroupVersion.Version,
		Resource: ResourceNameServiceEntry,
	},
	schema.GroupVersionResource{
		Group:    v1alpha3.GroupName,
		Version:  v1alpha3.SchemeGroupVersion.Version,
		Resource: ResourceNameEnvoyFilter,
	},
}
