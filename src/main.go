package main

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"PieFileMigrate/src/core"
	"PieFileMigrate/src/ui"
	"fmt"
)

func main() {
	fmt.Printf(constant.ConsolePrintGreen, " * 程序启动中 ")
	base.InitBase()
	core.InitCore()
	fmt.Printf(constant.ConsolePrintGreen, " * 程序启动成功 ")
	ui.InitHttpRouter()
}
