package base

import (
	"PieFileMigrate/src/constant"
	"PieFileMigrate/src/util"
	"gopkg.in/yaml.v3"
)

var ApplicationConfig applicationConfigModel
var HttpConfig httpConfigModel
var FtpConfig ftpConfigModel
var S3Config s3ConfigModel
var BoltDBConfig boltDBConfigModel
var RedisConfig redisConfigModel

// 应用配置模版
type applicationConfigModel struct {
	Application struct {
		Workers      []WorkerConfigModel `yaml:"workers"`
		Mq           MqConfigModel       `yaml:"mq"`
		ServerPort   int                 `yaml:"server-port"`
		MigrateMode  string              `yaml:"migrate-mode"`
		StorageMedia string              `yaml:"storage-media"`
		Log          Log                 `yaml:"log"`
	} `yaml:"application"`
}

type MqConfigModel struct {
	MaxSize      int `yaml:"max-size"`
	ConsumeBatch int `yaml:"consume-batch"`
	ConsumerNum  int `yaml:"consumer-num"`
}

type WorkerConfigModel struct {
	SourcePath          string `yaml:"source-path"`
	MigrateFileAgeLimit int64  `yaml:"migrate-file-age-limit"`
	JobCron             string `yaml:"job-cron"`
}

type Log struct {
	MaxAge int `yaml:"maxAge"`
}

// HTTP配置模版
type httpConfigModel struct {
	Http struct {
		TargetUrl string `yaml:"target-url"`
		Token     string `yaml:"token"`
	} `yaml:"http"`
}

// FTP配置模版
type ftpConfigModel struct {
	Ftp struct {
		TargetAddr string `yaml:"target-addr"`
		TargetPath string `yaml:"target-path"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
	} `yaml:"ftp"`
}

// S3配置模版
type s3ConfigModel struct {
	S3 struct {
		Endpoint  string `yaml:"endpoint"`
		Region    string `yaml:"region"`
		Bucket    string `yaml:"bucket"`
		AccessKey string `yaml:"access-key"`
		SecretKey string `yaml:"secret-key"`
		PathStyle bool   `yaml:"path-style"`
	} `yaml:"s3"`
}

// BoltDB数据库配置模版
type boltDBConfigModel struct {
	Boltdb struct {
		Db        string `yaml:"db"`
		TableName string `yaml:"table-name"`
	} `yaml:"boltdb"`
}

// Redis数据库配置模版
type redisConfigModel struct {
	Redis struct {
		Addr        string `yaml:"addr"`
		Password    string `yaml:"password"`
		Db          int    `yaml:"db"`
		KeyPrefix   string `yaml:"key-prefix"`
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

// InitHttpConfig 加载HTTP配置
func InitHttpConfig() {
	httpConfigBytes := loadConfigFile("./config/http.yaml")
	err := yaml.Unmarshal(httpConfigBytes, &HttpConfig)
	if err != nil {
		LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	LogHandler.Println(constant.LogInfoTag, "HTTP文件上传配置信息加载成功", string(httpConfigBytes))
}

// InitFtpConfig 加载FTP配置
func InitFtpConfig() {
	ftpConfigBytes := loadConfigFile("./config/ftp.yaml")
	err := yaml.Unmarshal(ftpConfigBytes, &FtpConfig)
	if err != nil {
		LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	LogHandler.Println(constant.LogInfoTag, "FTP文件上传配置信息加载成功", string(ftpConfigBytes))
}

// InitS3Config 加载S3配置
func InitS3Config() {
	s3ConfigBytes := loadConfigFile("./config/s3.yaml")
	err := yaml.Unmarshal(s3ConfigBytes, &S3Config)
	if err != nil {
		LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	LogHandler.Println(constant.LogInfoTag, "S3文件上传配置信息加载成功", string(s3ConfigBytes))
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
