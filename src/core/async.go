package core

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
)

// 内置消息队列
var mq = make(chan messageModel, base.Config.Server.MQMaxSize)

// 消息模版
type messageModel struct {
	// 文件命名
	FileName string
	// 文件路径
	LocalFilePath string
}

// 初始化MQ
func initMQ() {
	initMQConsumer()
	base.LogHandler.Println(constant.LogInfoTag, "内置消息队列启动成功")
}

// 异步迁移文件
func asyncMigrateFile(fileName string, localFilePath string) {
	msg := messageModel{
		FileName:      fileName,
		LocalFilePath: localFilePath,
	}
	mq <- msg
}

// 初始化MQ消费者
func initMQConsumer() {
	go func() {
		for {
			consumeMessage()
		}
	}()
}

// 消费消息
func consumeMessage() {
	defer func() {
		err := recover()
		if err != nil {
			base.LogHandler.Println(constant.LogErrorTag, err)
		}
	}()
	msg := <-mq
	err := uploadHandler.UploadFile(msg.FileName, msg.LocalFilePath)
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, "文件上传失败", err)
		return
	}
	// 如果上传成功，将文件标记为已上传
	ok := storageHandler.MarkFile(msg.LocalFilePath)
	if !ok {
		base.LogHandler.Println(constant.LogErrorTag, "文件标记失败")
		return
	}
}
