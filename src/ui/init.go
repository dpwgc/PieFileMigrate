package ui

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var api = NewHttpApi()

// InitHttpRouter 初始化HTTP路由
func InitHttpRouter() {

	base.LogHandler.Println(constant.LogInfoTag, "加载HTTP服务路由")

	r := httprouter.New()
	r.GET("/worker/monitor", api.WorkerMonitor)

	err := http.ListenAndServe(fmt.Sprintf(":%v", base.ApplicationConfig.Application.ServerPort), r)
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
		return
	}
}
