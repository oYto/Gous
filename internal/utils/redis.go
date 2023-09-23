package utils

import (
	"Gous/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	redisConn *redis.Client
	redisOnce sync.Once
)

// 连接 redis
func initRedis() {
	redisConfig := config.GetGlobalConf().RedisConfig
	log.Infof("redisConfig ====== %+v", redisConfig)
	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	redisConn = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisConfig.PassWord,
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSile,
	})
	if redisConn == nil {
		panic("failed to call redis.NewClient")
	}
	// 连接测试以确保与 Redis 服务器的通信正常。
	res, err := redisConn.Set(context.Background(), "abc", 100, 60).Result()
	log.Infof("res=======%v,err======%v", res, err)
	_, err = redisConn.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to ping redis, err:%s")
	}
}

// CloseRedis 关闭 redis 连接
func CloseRedis() {
	redisConn.Close()
}

func GetRedisCLi() *redis.Client {
	redisOnce.Do(initRedis)
	return redisConn
}
