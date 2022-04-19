module github.com/shuxnhs/istio-dashboard

go 1.16

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/envoyproxy/go-control-plane v0.10.2-0.20220413133113-27659a1a988e
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.7
	github.com/golang/protobuf v1.5.2
	github.com/kiali/kiali v1.49.0
	github.com/spf13/viper v1.11.0
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	gorm.io/driver/mysql v1.3.3
	gorm.io/gorm v1.23.4
	istio.io/client-go v1.13.2
	istio.io/istio v0.0.0-20220415183222-f611f67505bb
	istio.io/pkg v0.0.0-20220413132305-0219672e2d79
	k8s.io/api v0.23.5
	k8s.io/apiextensions-apiserver v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	sigs.k8s.io/gateway-api v0.4.2
)
