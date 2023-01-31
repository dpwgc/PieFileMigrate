package base

import (
	"PieFileMigrate/src/constant"
	"PieFileMigrate/src/util"
	"gopkg.in/yaml.v3"
)

var ApplicationConfig applicationConfigModel
var BoltDBConfig boltDBConfigModel
var RedisConfig redisConfigModel

// 应用配置模版
type applicationConfigModel struct {
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

// BoltDB数据库配置模版
type boltDBConfigModel struct {
	Boltdb struct {
		Db string `yaml:"db"`
	} `yaml:"boltdb"`
}

// Redis数据库配置模版
type redisConfigModel struct {
	Redis struct {
		Addr        string `yaml:"addr"`
		Password    string `yaml:"password"`
		Db          int    `yaml:"db"`
		PoolSize    int    `yaml:"pool-size"`
		MinIdleConn int    `yaml:"min-idle-conn"`
		MaxConnAge  int    `yaml:"max-conn-age"`
	} `yaml:"redis"`
}

// 加载应用配置
func initApplicationConfig() {
	applicationConfigBytes := loadConfigFile("./config/application.yaml")
	err := yaml.Unmarshal(applicationConfigBytes, &ApplicationConfig)
	if err != nil {
		LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	LogHandler.Println(constant.LogInfoTag, "应用配置信息加载成功", string(applicationConfigBytes))
}

// InitBoltDBConfig 加载boltdb配置
func InitBoltDBConfig() {
	boltdbConfigBytes := loadConfigFile("./config/boltdb.yaml")
	err := yaml.Unmarshal(boltdbConfigBytes, &BoltDBConfig)
	if err != nil {
		LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	LogHandler.Println(constant.LogInfoTag, "本地数据库(BoltDB)配置信息加载成功", string(boltdbConfigBytes))
}

// InitRedisConfig 加载redis配置
func InitRedisConfig() {
	redisConfigBytes := loadConfigFile("./config/redis.yaml")
	err := yaml.Unmarshal(redisConfigBytes, &RedisConfig)
	if err != nil {
		LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	LogHandler.Println(constant.LogInfoTag, "线上数据库(Redis)配置信息加载成功", string(redisConfigBytes))
}

// 读取本地配置文件
func loadConfigFile(path string) []byte {
	//加载本地配置
	configBytes, err := util.ReadFile(path)
	if err != nil {
		LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	return configBytes
}
