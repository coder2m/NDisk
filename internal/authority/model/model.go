package model

import (
	"errors"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	xgorm "github.com/coder2z/component/xinvoker/gorm"
	"github.com/coder2z/g-saber/xlog"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mainDB   *gorm.DB
	casbinDB *gormadapter.Adapter
)

func MainDB() *gorm.DB {
	if mainDB == nil {
		mainDB = xgorm.Invoker("main")
	}
	return mainDB
}

func CasbinDB() *gormadapter.Adapter {
	if casbinDB == nil {
		var err error
		casbinDB, err = gormadapter.NewAdapterByDBUseTableName(MainDB(), "sys", "casbin")
		if !errors.Is(err, nil) {
			xlog.Error("CasbinDB Connection failed",
				xlog.FieldErr(err),
				xlog.FieldComponentName("casbin"),
			)
		}
	}
	return casbinDB
}
