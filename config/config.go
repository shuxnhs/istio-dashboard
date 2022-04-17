package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config GlobalConfig

type GlobalConfig struct {
	WebConfig   `yaml:"WebConfig"`
	LogConfig   `yaml:"LogConfig"`
	MySQLConfig `yaml:"MySQLConfig"`
}

type WebConfig struct {
	ListenPort int `yaml:"ListenPort" env:"LISTEN_PORT" envDefault:"9655"`
}

type LogConfig struct {
	LogDir   string `yaml:"LogDir" env:"LOG_DIR" envDefault:"logs"`
	LogLevel string `yaml:"LogLevel" env:"LOG_LEVEL" envDefault:"debug"`
}

type MySQLConfig struct {
	DbHost string `yaml:"DbHost" env:"DB_HOST" envDefault:"localhost"`
	DbPort string `yaml:"DbPort" env:"DB_PORT" envDefault:"3306"`
	DbUser string `yaml:"DbUser" env:"DB_USER" envDefault:"root"`
	DbPass string `yaml:"DbPass" env:"DB_PASS" envDefault:""`
	DbName string `yaml:"DbName" env:"DB_NAME" envDefault:"Db"`
}

func InitializeConfig() *viper.Viper {
	config := "./config/config.yaml"
	// 生产环境可以通过设置环境变量来改变配置文件路径
	if configEnv := os.Getenv("ISTIO_DASHBOARD_CONFIG"); configEnv != "" {
		config = configEnv
	}

	// 初始化 viper
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s \n", err))
	}

	// 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		// 重载配置
		if err := v.Unmarshal(&Config); err != nil {
			fmt.Println(err)
		}
	})
	// 将配置赋值给全局变量
	if err := v.Unmarshal(&Config); err != nil {
		fmt.Println(err)
	}

	return v
}
