package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Input(r *bufio.Reader) string {
	v, _, e := r.ReadLine()

	s := string(v)

	if e != nil {
		fmt.Println("输入错误，退出程序！")
		os.Exit(0)
	}

	if strings.EqualFold(s, "Q") {
		fmt.Println("您选择了 “退出程序”")
		os.Exit(0)
	}

	return s
}
