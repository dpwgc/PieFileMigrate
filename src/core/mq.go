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
	Flag bool
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
		Flag:     true,
		FileName: fileName,
		FilePath: filePath,
		ModTime:  modTime,
	}
	mq <- msg
}

// 初始化MQ消费者
func initMqConsumer() {
	for i := 0; i < base.ApplicationConfig.Application.Mq.ConsumerNum; i++ {
		if base.ApplicationConfig.Application.Mq.ConsumeBatch > 1 {
			go enableBatchConsumer()
		} else {
			go enableOneConsumer()
		}
		base.LogHandler.Println(constant.LogInfoTag, "内置消息队列", fmt.Sprintf("%v号消费者启动", i))
	}
}

// 启动批量消费者
func enableBatchConsumer() {
	// 队列长度低于一定阈值后，后台协程每隔一秒往队列里塞个空消息，避免消费到最后一组时，消息数量不够导致无法正常推送
	go func() {
		for {
			if len(mq) < base.ApplicationConfig.Application.Mq.ConsumeBatch*base.ApplicationConfig.Application.Mq.ConsumerNum {
				msg := messageModel{
					Flag: false,
				}
				mq <- msg
			}
			time.Sleep(time.Second * 1)
		}
	}()
	// 循环批量消费消息
	for {
		batchConsumeMessage(1000)
	}
}

// 启动消费者
func enableOneConsumer() {
	// 循环消费消息
	for {
		oneConsumeMessage()
	}
}

// 批量消费消息
func batchConsumeMessage(batchSize int) {
	defer func() {
		err := recover()
		if err != nil {
			base.LogHandler.Println(constant.LogErrorTag, err)
		}
	}()
	var fileNames []string
	var filePaths []string
	var modTimes []time.Time
	for i := 0; i < batchSize; i++ {
		msg := <-mq
		//如果是空消息，跳过
		if !msg.Flag {
			continue
		}
		fileNames = append(fileNames, msg.FileName)
		filePaths = append(filePaths, msg.FilePath)
		modTimes = append(modTimes, msg.ModTime)
	}
	err := uploadHandler.UploadFiles(fileNames, filePaths, modTimes)
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, "文件上传失败", err)
		return
	}
	if len(fileNames) == 0 {
		return
	}
	// 如果上传成功，将文件标记为已上传
	for i := 0; i < len(fileNames); i++ {
		ok := storageHandler.MarkFile(filePaths[i], modTimes[i])
		if !ok {
			base.LogHandler.Println(constant.LogErrorTag, "文件标记失败")
			return
		}
	}
}

// 消费单条消息
func oneConsumeMessage() {
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
