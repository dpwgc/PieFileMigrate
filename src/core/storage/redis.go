package storage

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
)

func NewRedisStorageHandler() Handler {
	base.LogHandler.Println(constant.LogInfoTag, "线上数据库(Redis)加载成功")
	return &RedisStorageHandler{}
}

type RedisStorageHandler struct{}

func (s *RedisStorageHandler) MarkFile(filePath string) bool {
	return false
}

func (s *RedisStorageHandler) CheckFile(filePath string) bool {
	return false
}
