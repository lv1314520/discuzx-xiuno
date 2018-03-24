package xn3ToXn4

import (
	"bufio"
	"fmt"
	"github.com/skiy/xiuno-tools/lib"
	"log"
	"os"
)

type App struct {
}

type dbstr struct {
	lib.Database
	DBPre string
}

func (this *App) Init() {
	fmt.Println("\r\n===您选择了“1. Xiuno3 升级到 Xiuno4”===\r\n")

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

	xn3db, err := db3str.Connect()
	if err != nil {
		fmt.Println(err)
		log.Fatalln("\r\nXiuno3 数据库配置错误")
	}
	defer xn3db.Close()

	err = xn3db.Ping()
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

	xn4db, err := db4str.Connect()
	if err != nil {
		log.Fatalln("\r\nXiuno4 数据库配置错误")
	}
	defer xn4db.Close()

	err = xn4db.Ping()
	if err != nil {
		log.Fatalln("\r\nXiuno4: " + err.Error())
	}

	if db3str.DBHost == db4str.DBHost && db3str.DBName == db4str.DBName {
		log.Fatalln("\r\n不能在同一个数据库里升级，否则数据会被清空！请将新论坛安装到其他数据库。")
	}

	tables := [...]string{"group", "user"}
	for i, table := range tables {
		fmt.Println(string(i) + "正在转换表: " + table)

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
		}
	}
}
