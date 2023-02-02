package ui

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func NewHttpApi() Api {
	return &HttpApi{}
}

type Api interface {
	WorkerMonitor(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

type HttpApi struct{}

// WorkerMonitor 迁移工作者监控接口
func (a *HttpApi) WorkerMonitor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	m := make(map[string]base.WorkerMonitorModel)
	base.WorkerMonitorMap.Range(func(k, v interface{}) bool {
		m[k.(string)] = v.(base.WorkerMonitorModel)
		return true
	})
	jsonStr, err := json.Marshal(m)
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, err)
		return
	}
	_, err = w.Write(jsonStr)
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, err)
		return
	}
}
