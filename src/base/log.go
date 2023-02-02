package base

import (
	"log"
	"os"
	"time"
)

var LogHandler *log.Logger

func initLog() {
	file := "./log/runtime." + time.Now().Format("2006-01-01") + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		//创建log目录
		err = os.Mkdir("./log", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	LogHandler = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}
