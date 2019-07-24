package main

import (
	"fmt"
	"github.com/skiy/xiuno-tools/app"
	"github.com/skiy/xiuno-tools/app/libraries/database"
	"github.com/skiy/xiuno-tools/app/libraries/mcfg"
	"github.com/skiy/xiuno-tools/app/libraries/mlog"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Printf(`
:::
::: Discuz!X 3.x 转 XiunoBBS 4.x 工具
:::
::: 作者: Skiychan <dev@skiy.net> https://www.skiy.net
::: 本项目讨论帖：https://bbs.jadehive.com/thread-8059.htm
:::
::: 执行过程中按 "Ctrl + Z" 结束本程序...
:::
::: Version: 2.0.1    Updated: 2019-07-24
:::

`)

	//配置初始化
	initialize()

	// 判断 MYSQL 连接是否正常
	if err := checkConnectDB(); err != nil {
		mlog.Log.Fatal("", "数据库连接失败: %s", err.Error())
	}

	start := time.Now()

	mlog.Log.Info("", "开始导入数据 ...")

	app.NewApp().Parsing()

	mlog.Log.Info("", "已将 Discuz!X 转换至 XiunoBBS, 总耗时: %v\n", time.Since(start))

	fmt.Printf(`
:::
::: 本项目开源地址: https://github.com/skiy/xiuno-tools
::: 开发者 QQ: 869990770 技术支持论坛: https://bbs.jadehive.com
::: 如需技术支持请加 QQ 群: 891844359
:::
::: 如有意见和建议或者遇到 BUG，请到 GitHub 提 issue 。
:::

`)

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
		return fmt.Errorf("%s(Discuz!X)", err.Error())
	}

	if err = database.GetUcDB().PingMaster(); err != nil {
		return fmt.Errorf("%s(Discuz!UCenter)", err.Error())
	}

	if err = database.GetXiunoDB().PingMaster(); err != nil {
		return fmt.Errorf("%s(XiunoBBS)", err.Error())
	}

	return
}
