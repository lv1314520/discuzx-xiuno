package mlog

import (
	"fmt"
	"github.com/gogf/gf/g/os/glog"
	"xiuno-tools/app/libraries/mcfg"
)

var (
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
		"warn":          glog.LEVEL_WARN | glog.LEVEL_ERRO,
		"warning":       glog.LEVEL_WARN | glog.LEVEL_ERRO,
		"error":         glog.LEVEL_ERRO,
		"critical":      glog.LEVEL_CRIT,
	}
)

func InitLog() *Mlogger {
	Log = NewLogger()
	return Log
}

func ReadLog() *Mlogger {
	cfg := mcfg.GetCfg()

	//日志等级
	if logLevel := cfg.GetString("setting.log_level"); logLevel != "" {
		if l, ok := level[logLevel]; ok {
			Log.SetLevel(l)
		}
	}

	//日志路径
	if logPath := cfg.GetString("setting.log_path"); logPath != "" {
		Log.Logger.Path(logPath)
	}

	return Log
}

func Debugln(flag string, v ...interface{}) {
	Log.Debug(flag, fmt.Sprintln(v...))
}

func Fatalln(flag string, v ...interface{}) {
	Log.Fatal(flag, fmt.Sprintln(v...))
}
