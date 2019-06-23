package main

import (
	"errors"
	"fmt"
	"runtime"
	"time"
	"xiuno-tools/app"
	"xiuno-tools/app/libraries/database"
	"xiuno-tools/app/libraries/mcfg"
	"xiuno-tools/app/libraries/mlog"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println(`
:::
::: 本项目开源地址: https://github.com/skiy/xiuno-tools
::: 作者: Skiychan <dev@skiy.net> https://www.skiy.net
::: 本项目讨论帖：https://bbs.jadehive.com/thread-8059.htm
:::
:::
::: 执行过程中按 "Ctrl + Z" 结束本程序...
:::
::: Version:2.0.0    Updated:2019-05-01
:::
`)

	//配置初始化
	initialize()

	// 判断 MYSQL 连接是否正常
	if err := checkConnectDB(); err != nil {
		mlog.Log.Fatal("", "数据库连接失败: %s", err.Error())
	}

	start := time.Now()

	app.NewApp().Parsing()

	fmt.Println(fmt.Sprintf(`
:::
::: 已将 Discuz!X 升级至 XiunoBBS
::: 总耗时: %v
:::
::: 本程序开源地址: https://github.com/skiy/xiuno-tools
::: 作者: Skiychan <dev@skiy.net> https://www.skiy.net
::: QQ: 86999077 技术支持论坛: https://bbs.jadehive.com
:::
::: 如有意见和建议或者遇到 BUG，请到 GitHub 提 issue 。
:::`, time.Since(start)))

}

/**
配置初始化
*/
func initialize() {

	//配置文件
	mcfg.InitCfg()

	//日志跟踪
	mlog.InitTrace()

	//日志初始化
	mlog.InitLog()

	//日志配置
	mlog.ReadLog()

	//初始化前缀
	database.InitPrefix()
}

/**
检测数据库连接是否正常
*/
func checkConnectDB() (err error) {
	if err = database.GetDiscuzDB().PingMaster(); err != nil {
		return errors.New(fmt.Sprintf("%s(Discuz!X)", err.Error()))
	}

	if err = database.GetUcDB().PingMaster(); err != nil {
		return errors.New(fmt.Sprintf("%s(Discuz!UCenter)", err.Error()))
	}

	if err = database.GetXiunoDB().PingMaster(); err != nil {
		return errors.New(fmt.Sprintf("%s(XiunoBBS)", err.Error()))
	}

	return
}

/**
检测配置文件参数是否完整
*/
func checkCfg() {

}
