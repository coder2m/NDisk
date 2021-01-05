package model

import (
	xgorm "github.com/myxy99/component/xinvoker/gorm"
	"gorm.io/gorm"
)

var (
	mainDB *gorm.DB
)

func MainDB() *gorm.DB {
	if mainDB == nil {
		mainDB = xgorm.Invoker("main")
	}
	return mainDB
}
