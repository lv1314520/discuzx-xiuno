package app

import (
	"errors"
	"fmt"
	"xiuno-tools/app/controllers"
	"xiuno-tools/app/libraries/mcfg"
	"xiuno-tools/app/libraries/mlog"
)

type app struct {
}

func NewApp() *app {
	t := &app{}
	return t
}

func (t *app) Parsing() {
	tablesName := [...]string{
		"user",
		"group",
		"forum",
		"attach",
		"thread",
		"post",
		"thread_top",
		"mythread",
		"mypost",
	}

	var err error
	var ctrl controllers.Controller

	//mlog.Debugln("", mcfg.GetCfg().GetArray("tables.xiuno"))

	// 遍历表
	for _, table := range tablesName {
		cfgOffset := fmt.Sprintf("tables.xiuno.%s", table)

		// 是否转换
		if mcfg.GetCfg().GetBool(fmt.Sprintf("%s.convert", cfgOffset)) {
			// 转换控制器
			if ctrl, err = t.ctrl(table); err != nil {
				mlog.Log.Warning("", "%s(%s)", err.Error(), table)
				continue
			}

			if err = ctrl.ToConvert(); err != nil {
				mlog.Log.Fatal("", "转换数据表(%s)失败: %s", table, err.Error())
			}
		}
	}
}

/**
返回控制器
*/
func (t *app) ctrl(name string) (ctrl controllers.Controller, err error) {
	switch name {
	case "user":
		ctrl = controllers.NewUser()

	case "group":
		ctrl = controllers.NewGroup()

	case "forum":
		ctrl = controllers.NewForum()

	case "attach":
		ctrl = controllers.NewAttach()

	case "thread":
		ctrl = controllers.NewThread()

	case "post":
		ctrl = controllers.NewPost()

	case "thread_top":
		ctrl = controllers.NewThreadTop()

	case "mythread":
		ctrl = controllers.NewMythread()

	case "mypost":
		ctrl = controllers.NewMypost()

	default:
		err = errors.New("找不到对应的控制器,无法转换该表")

	}

	return
}
