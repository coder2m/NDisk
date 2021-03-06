package model

import (
	xgorm "github.com/coder2z/component/xinvoker/gorm"
	xredis "github.com/coder2z/component/xinvoker/redis"
	"github.com/go-redis/redis/v8"
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
