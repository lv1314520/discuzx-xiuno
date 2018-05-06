package main

import (
	"bufio"
	"fmt"
	"github.com/skiy/golib"
	"github.com/skiy/xiuno-tools/dx3ToXn4"
	"github.com/skiy/xiuno-tools/xn3ToXn4"
	"os"
)

func main() {
	fmt.Println(`
:::
::: 本程序开源地址: https://github.com/skiy/xiuno-tools
::: 作者: Skiychan <dev@skiy.net> https://www.skiy.net
:::
::: 请选择主菜单:::
:::
::: 1. XiunoBBS 3.x 升级到 XiunoBBS 4.x
::: 2. Discuz!X 3.x 升级到 XiunoBBS 4.x
:::
::: 执行过程中按"Q", 再按"回车键"退出本程序...
:::
::: Version:1.1.3    Updated:2018-05-06
`)

	var flag bool
	buf := bufio.NewReader(os.Stdin)
	for {
		inputVal := lib.Input(buf)

		switch inputVal {
		case "1":
			flag = true
			app := xn3ToXn4.App{}
			app.Init()
			break

		case "2":
			flag = true
			app := dx3ToXn4.App{}
			app.Init()
			break

		case "q":
		case "Q":
			flag = true
			break
		}

		//Q退出
		if flag {
			break
		}

		fmt.Println(`
-----------------------------------------
输入错误，请按上面提示选择对应的菜单
       按“Q”退出本程序
-----------------------------------------`)
	}
}
