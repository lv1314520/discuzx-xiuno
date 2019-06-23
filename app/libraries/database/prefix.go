package database

import (
	"github.com/gogf/gf/g/container/gmap"
	"xiuno-tools/app/libraries/mcfg"
)

var (
	prefix *gmap.Map
)

func InitPrefix() *gmap.Map {
	prefix = gmap.New()

	cfg := mcfg.GetCfg()
	u := cfg.GetString("database.uc.0.prefix")
	d := cfg.GetString("database.discuz.0.prefix")
	x := cfg.GetString("database.xiuno.0.prefix")

	prefix.Set("uc", u)
	prefix.Set("discuz", d)
	prefix.Set("xiuno", x)

	return prefix
}

func AddPrefix(key string, pre string) {
	prefix.Set(key, pre)
}

func GetPrefix(key string) string {
	if res := prefix.Get(key); res == nil {
		return ""
	} else {
		return res.(string)
	}
}

func GetPrefixs() *gmap.Map {
	return prefix
}

func Remove(key string) {
	prefix.Remove(key)
}
