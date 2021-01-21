package model

import (
	"github.com/go-redis/redis/v8"
	xgorm "github.com/myxy99/component/xinvoker/gorm"
	xredis "github.com/myxy99/component/xinvoker/redis"
	"gorm.io/gorm"
)

var (
	mainDB    *gorm.DB
	mainRedis *redis.Client
)

func MainDB() *gorm.DB {
	if mainDB == nil {
		mainDB = xgorm.Invoker("main")
	}
	return mainDB
}

func MainRedis() *redis.Client {
	if mainRedis == nil {
		mainRedis = xredis.Invoker("main")
	}
	return mainRedis
}
