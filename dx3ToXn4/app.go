package dx3ToXn4

import (
	"fmt"
	"github.com/skiy/xiuno-tools/lib"
)

type App struct {
}

type dbstr struct {
	lib.Database
	DBPre string
	Auto  bool
}

func (this *App) Init() {
	fmt.Println("\r\n===您选择了: 2. Discuz!X3.x 升级到 Xiuno4\r\n")
}
