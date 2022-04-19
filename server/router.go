package server

import (
	"net/http"

	"github.com/shuxnhs/istio-dashboard/api"
	_ "github.com/shuxnhs/istio-dashboard/docs"
	"github.com/shuxnhs/istio-dashboard/server/middleware"

	"github.com/gin-gonic/gin"
	ginSwaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, Istio Dashboard")
	})

	r.GET("/healthy", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(ginSwaggerFiles.Handler))

	// 路由
	project := r.Group("/project")
	{
		project.GET("list", api.ListProjects)
	}

	kube := r.Group("/kube")
	{
		namespace := kube.Group("/namespace")
		{
			namespace.GET("list", api.ListNamespace)
		}

	}

	sidecar := r.Group("/sidecar")
	{
		sidecar.GET("check", api.Check)

		eds := sidecar.Group("/eds")
		{
			eds.GET("list", api.ListEDS)
		}

		cds := sidecar.Group("/cds")
		{
			cds.GET("list", api.ListCDS)
		}

	}
	return r
}
