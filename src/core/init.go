package core

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"PieFileMigrate/src/core/storage"
	"PieFileMigrate/src/core/upload"
	"fmt"
)

var storageHandler storage.Handler
var uploadHandler upload.Handler

// InitCore 加载核心服务
func InitCore() {

	//存储介质选择
	switch base.ApplicationConfig.Application.StorageMedia {
	case "boltdb":
		storageHandler = storage.NewBoltDBStorageHandler()
		break
	case "redis":
		storageHandler = storage.NewRedisStorageHandler()
		break
	}

	//迁移方式选择
	switch base.ApplicationConfig.Application.MigrateMode {
	case "http":
		uploadHandler = upload.NewHTTPUploadHandler()
		break
	case "ftp":
		uploadHandler = upload.NewFTPUploadHandler()
		break
	case "s3":
		uploadHandler = upload.NewS3UploadHandler()
	}

	//加载消息队列
	initMq()

	//启动定时任务
	initWorker()

	base.LogHandler.Println(constant.LogInfoTag, "核心服务加载成功")
	fmt.Printf(constant.ConsolePrintCyan, " * 核心服务加载成功 ")
}
