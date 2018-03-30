package xn3ToXn4

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/skiy/xiuno-tools/lib"
	"log"
	"os"
	"strings"
)

type App struct {
}

type dbstr struct {
	lib.Database
	DBPre string
	Auto  bool
}

var xiuno3db, xiuno4db *sql.DB

func (this *App) Init() {
	fmt.Println("\r\n===您选择了: 1. Xiuno3 升级到 Xiuno4\r\n")

	db3str := dbstr{}
	fmt.Println("正在配置 Xiuno3 数据库")
	db3str.Setting()

	buf := bufio.NewReader(os.Stdin)
	fmt.Println("请配置数据库表前缀:(默认为 bbs_)")
	s := lib.Input(buf)
	if s == "" {
		s = "bbs_"
	}
	db3str.DBPre = s
	fmt.Println("数据库表前缀为: " + s)

	var err error
	xiuno3db, err = db3str.Connect()
	if err != nil {
		fmt.Println(err)
		log.Fatalln("\r\nXiuno3 数据库配置错误")
	}

	err = xiuno3db.Ping()
	if err != nil {
		log.Fatalln("\r\nXiuno3: " + err.Error())
	}

	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	db4str := dbstr{}
	fmt.Println("正在配置 Xiuno4 数据库")
	db4str.Setting()

	buf = bufio.NewReader(os.Stdin)
	fmt.Println("请配置数据库表前缀:(默认为 bbs_)")
	s = lib.Input(buf)
	if s == "" {
		s = "bbs_"
	}
	db4str.DBPre = s
	fmt.Println("数据库表前缀为: " + s)

	xiuno4db, err = db4str.Connect()
	if err != nil {
		log.Fatalln("\r\nXiuno4 数据库配置错误")
	}

	err = xiuno4db.Ping()
	if err != nil {
		log.Fatalln("\r\nXiuno4: " + err.Error())
	}

	if db3str.DBHost == db4str.DBHost && db3str.DBName == db4str.DBName {
		log.Fatalln("\r\n不能在同一个数据库里升级，否则数据会被清空！请将新论坛安装到其他数据库。")
	}

	buf = bufio.NewReader(os.Stdin)
	fmt.Println("全自动更新所有表(Y/N): (默认为 Y)")
	s = lib.Input(buf)
	if !strings.EqualFold(s, "N") {
		db4str.Auto = true
	}

	tables := [...]string{
		"group",
		"user",
		"user_open_plat",
		"forum",
		"forum_access",
		"attach",
		"modlog",
		"friendlink",
		"thread",
		"post",
	}

	for _, table := range tables {
		fmt.Println("正在转换表: " + table)

		switch table {
		case "group":
			do := group{}
			do.db3str = db3str
			do.db4str = db4str

			do.update()
			break

		case "user":
			do := user{}
			do.db3str = db3str
			do.db4str = db4str

			do.update()
			break

		case "user_open_plat":
			do := user_open_plat{}
			do.db3str = db3str
			do.db4str = db4str

			do.update()
			break

		case "post":
			do := post{}
			do.db3str = db3str
			do.db4str = db4str

			do.update()
			break

		case "forum":
			do := forum{}
			do.db3str = db3str
			do.db4str = db4str

			do.update()
			break

		case "forum_access":
			do := forum_access{}
			do.db3str = db3str
			do.db4str = db4str

			do.update()
			break

		case "thread":
			do := thread{}
			do.db3str = db3str
			do.db4str = db4str

			do.update()
			break

		case "attach":
			do := attach{}
			do.db3str = db3str
			do.db4str = db4str

			do.update()
			break

		case "modlog":
			do := modlog{}
			do.db3str = db3str
			do.db4str = db4str

			do.update()
			break

		case "friendlink":
			do := friendlink{}
			do.db3str = db3str
			do.db4str = db4str

			do.update()
			break
		}
	}

	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(`
:::
::: 已将 Xinno3 升级至 Xiuno4
::: 您还需要将 xiuno3 下的 upload 移动到 xiuno4 下 
::: 本程序开源地址: https://github.com/skiy/xiuno-tools
::: 作者: Skiychan <dev@skiy.net> https://www.skiy.net
:::
`)
}
