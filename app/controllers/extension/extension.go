package extension

import (
	"github.com/gogf/gf/g/os/gcfg"
	"xiuno-tools/app/libraries/mcfg"
	"xiuno-tools/app/libraries/mlog"
)

var cfg *gcfg.Config

type extension struct {
}

func NewExtension() *extension {
	t := &extension{}
	return t
}

func (t *extension) Parsing() {
	cfg = mcfg.GetCfg()

	t.ShowError(NewGroup().Parsing())
	t.ShowError(NewUser().Parsing())
	t.ShowError(NewThreadPost().Parsing())
}

func (t *extension) ShowError(err error) {
	if err != nil {
		mlog.Log.Warning("", "%s", err.Error())
	}
}
