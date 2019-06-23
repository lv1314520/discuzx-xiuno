package mlog

import "github.com/gogf/gf/g/os/glog"

type Mlogger struct {
	Logger *glog.Logger
}

func NewLogger() *Mlogger {
	t := &Mlogger{}
	t.Logger = glog.New()
	t.Logger.SetLevel(glog.LEVEL_ALL)
	return t
}

func (t *Mlogger) SetLevel(l int) {
	t.Logger.SetLevel(l)
	//t.Logger.Backtrace(true)
}

func (t *Mlogger) Critical(flag, format string, v ...interface{}) {
	t.Logger.Criticalf(format, v...)
	if flag != "" {
		if tr := GetTrace(flag); tr != nil {
			tr.Log.Criticalf(format, v...)
		}
	}
}

func (t *Mlogger) Error(flag, format string, v ...interface{}) {
	t.Logger.Errorf(format, v...)
	if flag != "" {
		if tr := GetTrace(flag); tr != nil {
			tr.Log.Errorf(format, v...)
		}
	}
}

func (t *Mlogger) Notice(flag, format string, v ...interface{}) {
	t.Logger.Noticef(format, v...)
	if flag != "" {
		if tr := GetTrace(flag); tr != nil {
			tr.Log.Noticef(format, v...)
		}
	}
}

func (t *Mlogger) Debug(flag, format string, v ...interface{}) {
	t.Logger.Debugf(format, v...)
	if flag != "" {
		if tr := GetTrace(flag); tr != nil {
			tr.Log.Debugf(format, v...)
		}
	}
}

func (t *Mlogger) Warning(flag, format string, v ...interface{}) {
	t.Logger.Warningf(format, v...)
	if flag != "" {
		if tr := GetTrace(flag); tr != nil {
			tr.Log.Warningf(format, v...)
		}
	}
}

func (t *Mlogger) Info(flag, format string, v ...interface{}) {
	t.Logger.Infof(format, v...)
	if flag != "" {
		if tr := GetTrace(flag); tr != nil {
			tr.Log.Infof(format, v...)
		}
	}
}

/**
致命错误
*/
func (t *Mlogger) Fatal(flag, format string, v ...interface{}) {
	t.Logger.Fatalf(format, v...)
	if flag != "" {
		if tr := GetTrace(flag); tr != nil {
			tr.Log.Fatalf(format, v...)
		}
	}
}

func (t *Mlogger) Show(flag, format string, v ...interface{}) {
	if flag != "" {
		if tr := GetTrace(flag); tr != nil {
			tr.Log.Debugf(format, v...)
		}
	}
}
