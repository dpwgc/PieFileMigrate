package main

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"PieFileMigrate/src/core"
	"fmt"
	"time"
)

func main() {
	fmt.Printf(constant.ConsolePrintGreen, " * 程序启动中 ")
	base.InitBase()
	core.InitCore()
	fmt.Printf(constant.ConsolePrintGreen, " * 程序启动成功 ")
	for {
		time.Sleep(100 * time.Second)
	}
}
