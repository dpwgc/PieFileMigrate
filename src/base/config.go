package base

import (
	"PieFileMigrate/src/constant"
	"PieFileMigrate/src/util"
	"gopkg.in/yaml.v3"
)

var Config configModel

// 配置模版
type configModel struct {
	Server struct {
		SourcePath           string `yaml:"source-path"`
		MigrateFileTimeLimit int64  `yaml:"migrate-file-time-limit"`
		JobCron              string `yaml:"job-cron"`
		TargetAddr           string `yaml:"target-addr"`
		MQMaxSize            int64  `yaml:"mq-max-size"`
		MigrateMode          string `yaml:"migrate-mode"`
		StorageMedia         string `yaml:"storage-media"`
	} `yaml:"server"`
}

// 本地配置加载
func initConfig() {
	//加载客户端配置
	configBytes, err := util.ReadFile("./config.yaml")
	if err != nil {
		LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	err = yaml.Unmarshal(configBytes, &Config)
	if err != nil {
		LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	LogHandler.Println(constant.LogInfoTag, "本地配置加载成功", string(configBytes))
}
