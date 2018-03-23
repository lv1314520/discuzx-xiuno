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
		log.Fatalln("Xiuno3 数据库配置错误")
	}

	err = xn3db.Ping()
	if err != nil {
		log.Fatalln("\r\nXiuno3: " + err.Error())
	}

	fmt.Println("")

	db4 := lib.Database{}
	fmt.Println("正在配置 Xiuno4 数据库")
	db4.Setting()

	xn4db, err := db4.Connect()
	if err != nil {
		log.Fatalln("Xiuno4 数据库配置错误")
	}

	err = xn4db.Ping()
	if err != nil {
		log.Fatalln("\r\nXiuno4: " + err.Error())
	}

	fmt.Println(xn3db.Ping(), xn4db.Ping())
}
