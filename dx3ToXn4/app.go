package dx3ToXn4

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/skiy/golib"
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

var (
	dxdb, xndb *sql.DB
	groupReset bool
)

func (this *App) Init() {
	oldname := "Discuz!X3.x"
	newname := "XiunoBBS4.x"

	fmt.Printf("\r\n===您选择了: 2. %s 升级到 %s\r\n\r\n", oldname, newname)

	dxstr := dbstr{}
	fmt.Printf("正在配置 %s 数据库\r\n", oldname)
	dxstr.Setting()

	buf := bufio.NewReader(os.Stdin)
	fmt.Println("请配置数据库表前缀:(空格为无前缀, 默认为 pre_)")
	s := lib.Input(buf)
	if s == "" {
		s = "pre_"
	} else if s == " " {
		s = ""
	}
	dxstr.DBPre = s
	fmt.Printf("数据库表前缀为: %s\r\n\r\n", s)

	var err error
	dxdb, err = dxstr.Connect()
	if err != nil {
		fmt.Println(err)
		log.Fatalf("\r\n 数据库配置错误\r\n", oldname)
	}

	err = dxdb.Ping()
	if err != nil {
		log.Fatalf("\r\n %s: %s\r\n", oldname, err.Error())
	}

	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	xnstr := dbstr{}
	fmt.Printf("正在配置 %s 数据库\r\n", newname)
	xnstr.Setting()

	fmt.Println("请配置数据库表前缀:(空格为无前缀, 默认为 bbs_)")
	s = lib.Input(buf)
	if s == "" {
		s = "bbs_"
	} else if s == " " {
		s = ""
	}

	xnstr.DBPre = s
	fmt.Printf("数据库表前缀为: %s\r\n\r\n", s)

	xndb, err = xnstr.Connect()
	if err != nil {
		fmt.Println(err)
		log.Fatalf("\r\n 数据库配置错误\r\n", newname)
	}

	err = xndb.Ping()
	if err != nil {
		log.Fatalf("\r\n %s: %s\r\n", newname, err.Error())
	}

	if dxstr.DBHost == xnstr.DBHost && dxstr.DBName == xnstr.DBName {
		if dxstr.DBPre == xnstr.DBPre {
			log.Fatalln(`
Discuz!X和XiunoBBS的数据库名和表前缀不能相同，
否则数据可能会造成破坏！
请将新论坛安装到其他数据库。
`)
		}

		fmt.Printf(`
您将在同一个数据库转换数据，
需要再次确认是否强行转换(Y/N): (默认为 N)
database host: %s 
database name: %s
`, xnstr.DBHost, xnstr.DBName)

		s = lib.Input(buf)
		if !strings.EqualFold(s, "Y") {
			log.Fatalln(`您已退出程序`)
		}
	}

	dxdb.SetMaxIdleConns(0)
	xndb.SetMaxIdleConns(0)

	fmt.Println("\r\n全自动更新所有表(Y/N): (默认为 Y)")
	s = lib.Input(buf)
	if !strings.EqualFold(s, "N") {
		xnstr.Auto = true
	}
	fmt.Println("")

	tables := [...]string{
		"user",
		"group",
		"forum",
		"attach",
		"thread",
		"post",
	}

	for _, table := range tables {
		fmt.Println("正在转换表: " + table)

		switch table {

		case "user":
			do := &user{}
			do.dxstr = dxstr
			do.xnstr = xnstr

			do.update()
			break

		case "group":
			do := &group{}
			do.dxstr = dxstr
			do.xnstr = xnstr

			do.update()
			break

		case "forum":
			do := &forum{}
			do.dxstr = dxstr
			do.xnstr = xnstr

			do.update()
			break

		case "thread":
			do := &thread{}
			do.dxstr = dxstr
			do.xnstr = xnstr

			do.update()
			break

		case "post":
			do := &post{}
			do.dxstr = dxstr
			do.xnstr = xnstr

			do.update()
			break

		case "attach":
			do := &attach{}
			do.dxstr = dxstr
			do.xnstr = xnstr

			do.update()
			break
		}
	}

	do := &extension{}
	do.dxstr = dxstr
	do.xnstr = xnstr

	do.update()

	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Println(`
:::
::: 已将 ` + oldname + ` 升级至 ` + newname + `
::: 本程序开源地址: https://github.com/skiy/xiuno-tools
::: 作者: Skiychan <dev@skiy.net> https://www.skiy.net
::: 如有意见和建议或者遇到 BUG，请到 GitHub 提 issue 。
:::
`)
}
