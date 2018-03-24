package main

import (
	"bufio"
	"fmt"
	"github.com/skiy/xiuno-tools/lib"
	"github.com/skiy/xiuno-tools/xn3ToXn4"
	"os"
)

func main() {
	buf := bufio.NewReader(os.Stdin)

	fmt.Println(`
::: 请选择主菜单...
:::
::: 1. Xiuno 3.x 升级到 Xiuno 4.x
:::
::: 执行过程中按"Q", 再按"回车键"退出本程序...
:::
::: 本程序开源地址: https://github.com/skiy/xiuno-tools
::: 作者: Skiychan <dev@skiy.net> https://www.skiy.net
`)

	inputVal := lib.Input(buf)

	if inputVal == "1" {
		app := xn3ToXn4.App{}
		app.Init()
	}
}
