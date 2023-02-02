package core

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"fmt"
	"github.com/robfig/cron/v3"
)

// 初始化迁移工作者
func initWorker() {
	for _, v := range base.ApplicationConfig.Application.Workers {
		enableJob(v.JobCron, v.SourcePath, v.MigrateFileAgeLimit)
		base.LogHandler.Println(constant.LogInfoTag, fmt.Sprintf("[ %s ] 目录迁移工作者启动", v.SourcePath))
	}
}

// 启动任务
func enableJob(jobCron string, sourcePath string, migrateFileAgeLimit int64) {
	job := newWithSeconds()
	_, err := job.AddFunc(jobCron, func() {
		doMigrate(sourcePath, migrateFileAgeLimit)
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
