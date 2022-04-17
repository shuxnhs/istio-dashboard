package server

import (
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

	// 路由
	project := r.Group("/project")
	{
		project.GET("list", api.ListProjects)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(ginSwaggerFiles.Handler))
	return r
}
