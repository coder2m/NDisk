package model

import (
	"errors"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"github.com/myxy99/component/xcfg"
	xgorm "github.com/myxy99/component/xinvoker/gorm"
	"github.com/myxy99/component/xlog"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mainDB   *gorm.DB
	casbinDB *xormadapter.Adapter
)

func MainDB() *gorm.DB {
	if mainDB == nil {
		mainDB = xgorm.Invoker("main")
	}
	return mainDB
}

func CasbinDB() *xormadapter.Adapter {
	if casbinDB == nil {
		var err error
		casbinDB, err = xormadapter.NewAdapter(xcfg.GetString("casbin.driver"), xcfg.GetString("casbin.host"))
		if !errors.Is(err, nil) {
			xlog.Error("CasbinDB Connection failed",
				xlog.FieldErr(err),
				xlog.FieldComponentName("casbin"),
			)
		}
	}
	return casbinDB
}
