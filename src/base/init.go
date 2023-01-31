package base

import (
	"PieFileMigrate/src/constant"
	"fmt"
)

func InitBase() {
	initLog()
	initApplicationConfig()
	LogHandler.Println(constant.LogInfoTag, "基础服务加载成功")
	fmt.Printf(constant.ConsolePrintCyan, " * 基础服务加载成功 ")
}
