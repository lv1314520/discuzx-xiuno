package xn3ToXn4

import (
	"fmt"
	"github.com/skiy/xiuno-tools/lib"
	"log"
)

type App struct {
}

func (this *App) Init() {
	fmt.Println("\r\n===您选择了“1. Xiuno3 升级到 Xiuno4”===\r\n")

	db3 := lib.Database{}
	fmt.Println("正在配置 Xiuno3 数据库")
	db3.Setting()

	xn3db, err := db3.Connect()
	if err != nil {
		fmt.Println(err)
		log.Fatalln("\r\nXiuno3 数据库配置错误")
	}

	err = xn3db.Ping()
	if err != nil {
		log.Fatalln("\r\nXiuno3: " + err.Error())
	}

	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	db4 := lib.Database{}
	fmt.Println("正在配置 Xiuno4 数据库")
	db4.Setting()

	xn4db, err := db4.Connect()
	if err != nil {
		log.Fatalln("\r\nXiuno4 数据库配置错误")
	}

	err = xn4db.Ping()
	if err != nil {
		log.Fatalln("\r\nXiuno4: " + err.Error())
	}

	if db3.DBHost == db4.DBHost && db3.DBName == db4.DBName {
		log.Fatalln("\r\n不能在同一个数据库里升级，否则数据会被清空！请将新论坛安装到其他数据库。")
	}
}
