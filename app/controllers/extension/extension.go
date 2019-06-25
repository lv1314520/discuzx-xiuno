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

	var ctrl Controller

	ctrl = NewGroup()
	t.ShowError(ctrl.Parsing())

	ctrl = NewUser()
	t.ShowError(ctrl.Parsing())

	ctrl = NewThreadPost()
	t.ShowError(ctrl.Parsing())

	ctrl = NewForum()
	t.ShowError(ctrl.Parsing())
}

func (t *extension) ShowError(err error) {
	if err != nil {
		mlog.Log.Warning("", "%s", err.Error())
	}
}
