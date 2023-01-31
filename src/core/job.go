package core

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"github.com/robfig/cron/v3"
)

// 初始化定时任务
func initJob() {
	job := newWithSeconds()
	_, err := job.AddFunc(base.Config.Server.JobCron, func() {
		doMigrate()
	})
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	job.Start()
	base.LogHandler.Println(constant.LogInfoTag, "定时任务启动成功")
}

// 返回一个支持至 秒 级别的 cron
func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
