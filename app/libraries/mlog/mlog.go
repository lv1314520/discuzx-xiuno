package mlog

import (
	"fmt"
	"github.com/gogf/gf/g/os/glog"
	"xiuno-tools/app/libraries/mcfg"
)

var (
	// Log logger
	Log   *Mlogger
	level = map[string]int{
		"all":           glog.LEVEL_ALL,
		"dev":           glog.LEVEL_DEV,
		"development":   glog.LEVEL_DEV,
		"prod":          glog.LEVEL_PROD,
		"production":    glog.LEVEL_PROD,
		"debug":         glog.LEVEL_DEV,
		"info":          glog.LEVEL_INFO,
		"informational": glog.LEVEL_INFO,
		"notice":        glog.LEVEL_NOTI | glog.LEVEL_WARN | glog.LEVEL_ERRO,
		"warn":          glog.LEVEL_WARN,
		"warning":       glog.LEVEL_WARN,
		"error":         glog.LEVEL_ERRO,
		"critical":      glog.LEVEL_CRIT,
		"alert":         glog.LEVEL_PROD | glog.LEVEL_INFO,
	}
)

// InitLog log init
func InitLog() *Mlogger {
	Log = NewLogger()
	return Log
}

// ReadLog read log config
func ReadLog() *Mlogger {
	cfg := mcfg.GetCfg()

	//日志等级
	if logLevel := cfg.GetString("log.level"); logLevel != "" {
		if l, ok := level[logLevel]; ok {
			Log.SetLevel(l)
		}
	}

	//日志路径
	if logPath := cfg.GetString("log.path"); logPath != "" {
		if err := Log.Logger.SetPath(logPath); err != nil {
			Log.Logger.Warning(err.Error())
		}
	}

	// 是否输出错误行
	Log.Logger.SetBacktrace(cfg.GetBool("log.trace"))

	//Log.Logger.SetStdoutPrint(false)

	return Log
}

// Debugln debug
func Debugln(flag string, v ...interface{}) {
	Log.Debug(flag, fmt.Sprintln(v...))
}

// Fatalln fatal
func Fatalln(flag string, v ...interface{}) {
	Log.Fatal(flag, fmt.Sprintln(v...))
}
