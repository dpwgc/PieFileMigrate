package core

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"fmt"
	"time"
)

// 内置消息队列
var mq = make(chan messageModel, base.ApplicationConfig.Application.Mq.MaxSize)

// 消息模版
type messageModel struct {
	// 文件命名
	FileName string
	// 文件路径
	FilePath string
	// 文件更新时间
	ModTime time.Time
}

// 初始化内置消息队列
func initMq() {
	initMqConsumer()
	base.LogHandler.Println(constant.LogInfoTag, "内置消息队列启动成功")
}

// 异步迁移文件
func asyncMigrateFile(fileName string, filePath string, modTime time.Time) {
	msg := messageModel{
		FileName: fileName,
		FilePath: filePath,
		ModTime:  modTime,
	}
	mq <- msg
}

// 初始化MQ消费者
func initMqConsumer() {
	for i := 0; i < base.ApplicationConfig.Application.Mq.ConsumerNum; i++ {
		go enableConsumer()
		base.LogHandler.Println(constant.LogInfoTag, "内置消息队列", fmt.Sprintf("%v号消费者启动", i))
	}
}

// 启动消费者
func enableConsumer() {
	// 循环消费消息
	for {
		consumeMessage()
	}
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
	err := uploadHandler.UploadFile(msg.FileName, msg.FilePath, msg.ModTime)
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, "文件上传失败", err)
		return
	}
	// 如果上传成功，将文件标记为已上传
	ok := storageHandler.MarkFile(msg.FilePath, msg.ModTime)
	if !ok {
		base.LogHandler.Println(constant.LogErrorTag, "文件标记失败")
		return
	}
}
