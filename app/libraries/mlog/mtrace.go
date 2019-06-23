package mlog

import (
	"github.com/gogf/gf/g/container/gmap"
)

var (
	traces *gmap.Map
)

func InitTrace() {
	traces = gmap.New()
}

func AddTrace(key string, trace *LogTrace) {
	traces.Set(key, trace)
}

func GetTrace(key string) *LogTrace {
	if res := traces.Get(key); res == nil {
		return nil
	} else {
		return res.(*LogTrace)
	}
}

func GetTraces() *gmap.Map {
	return traces
}

func DeleteTrace(key string) {
	traces.Remove(key)
}
