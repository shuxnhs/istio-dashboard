package log

import (
	"time"

	"github.com/shuxnhs/istio-dashboard/config"

	"istio.io/pkg/log"
)

func InitializeLog() {
	date := time.Now().Format("2006-01-02")
	logOpts := log.DefaultOptions()
	logOpts.OutputPaths = append(logOpts.OutputPaths, "run.log")
	logOpts.RotateOutputPath = config.Config.LogDir + date + ".log"
	logOpts.WithStackdriverLoggingFormat()
	logOpts.SetOutputLevel("istio-dashboard", StringToLevel[config.Config.LogLevel])
	_ = log.Configure(logOpts)
}

// StringToLevel this is the same as istio.io/pkg/log.stringToLevel
var StringToLevel = map[string]log.Level{
	"debug": log.DebugLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"fatal": log.FatalLevel,
	"none":  log.NoneLevel,
}
