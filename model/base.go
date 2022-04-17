package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shuxnhs/istio-dashboard/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

var (
	KubeConfigDB *KubeConfig
)

// list add table name
const (
	// kubeConfig
	KubeConfigTableName = "kube_config"
)

// soft-delete
const (
	_             = iota
	StatusNormal  // 1:正常
	StatusDisable // 2:不可用
	StatusDeleted // 3:删除
)

func InitializeDatebase() {
	gormConfig := &gorm.Config{
		DisableAutomaticPing: false,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel:      logger.Info, // Log level
				Colorful:      false,       // 禁用彩色打印
			},
		),
	}

	dsn := GetDsn()
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		panic(err)
	}
}

// GetDsn 获取MySQL连接DSN
func GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&loc=Local&parseTime=true",
		config.Config.DbUser,
		config.Config.DbPass,
		config.Config.DbHost,
		config.Config.DbPort,
		config.Config.DbName)
}
