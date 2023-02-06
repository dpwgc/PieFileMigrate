package base

import (
	"PieFileMigrate/src/constant"
	"github.com/lestrrat-go/file-rotatelogs"
	"log"
	"time"
)

var LogHandler *log.Logger

func initLog() {
	writer, _ := rotatelogs.New(
		constant.LogFilePath,
		rotatelogs.WithMaxAge(time.Duration(ApplicationConfig.Application.Log.MaxAge*24)*time.Hour),
		rotatelogs.WithRotationTime(constant.LogFileRotationTime),
	)
	log.SetOutput(writer)
	LogHandler = log.Default()
}
