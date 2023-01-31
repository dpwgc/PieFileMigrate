package storage

import (
	"PieFileMigrate/src/base"
	"PieFileMigrate/src/constant"
	"PieFileMigrate/src/util"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewRedisStorageHandler() Handler {
	base.InitRedisConfig()
	client, err := initRedis()
	if err != nil {
		base.LogHandler.Println(constant.LogErrorTag, err)
		panic(err)
	}
	base.LogHandler.Println(constant.LogInfoTag, "线上数据库(Redis)加载成功")
	return &RedisStorageHandler{
		Client: client,
	}
}

type RedisStorageHandler struct {
	Client *redis.Client
}

func initRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         base.RedisConfig.Redis.Addr,
		Password:     base.RedisConfig.Redis.Password,
		DB:           base.RedisConfig.Redis.Db,
		PoolSize:     base.RedisConfig.Redis.PoolSize,
		MinIdleConns: base.RedisConfig.Redis.MinIdleConn,
		MaxConnAge:   time.Second * time.Duration(base.RedisConfig.Redis.MaxConnAge),
	})
	ping := client.Ping(context.Background())
	if ping.Err() != nil {
		return nil, ping.Err()
	}
	return client, nil
}

func (s *RedisStorageHandler) MarkFile(filePath string) bool {
	var ctx = context.Background()
	res := s.Client.Set(ctx, filePath, util.GetLocalDateTime(), -1)
	//写入数据库失败
	if res.Err() != nil {
		base.LogHandler.Println(constant.LogErrorTag, res.Err())
		return false
	}
	return true
}

func (s *RedisStorageHandler) CheckFile(filePath string) bool {
	var ctx = context.Background()
	res := s.Client.Get(ctx, filePath)
	//查询数据库失败
	if res.Err() != nil {
		base.LogHandler.Println(constant.LogErrorTag, res.Err())
		return false
	}
	return len(res.String()) > 0
}
