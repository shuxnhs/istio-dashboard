package main

import (
	"fmt"
	"github.com/shuxnhs/istio-dashboard/config"
	"github.com/shuxnhs/istio-dashboard/log"
	"github.com/shuxnhs/istio-dashboard/model"
	"github.com/shuxnhs/istio-dashboard/server"
)

/* ---This is for generate swagger--- */

// @title Istio-Dashboard
// @version 1.0
// @description Istio-Dashboard API

// @contact.name  sHuXnHs
// @contact.email 610087273@qq.com
// @BasePath /

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {
	config.InitializeConfig()
	log.InitializeLog()
	model.InitializeDatebase()

	// 装载路由
	r := server.NewRouter()
	r.Run(":" + fmt.Sprint(config.Config.ListenPort))
}
