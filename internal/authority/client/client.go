package xclient

import (
	"errors"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/coder2z/component/xcfg"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/ndisk/internal/authority/model"
)

var (
	casbinClient *casbin.Enforcer
)

func CasbinClient() *casbin.Enforcer {
	if casbinClient == nil {
		var (
			err error
			m   casbinModel.Model
		)
		m, err = casbinModel.NewModelFromString(xcfg.GetString("casbin.model"))
		if !errors.Is(err, nil) {
			xlog.Error("casbin Model failed",
				xlog.FieldErr(err),
				xlog.FieldComponentName("casbin"),
			)
			return nil
		}
		casbinClient, err = casbin.NewEnforcer(m, model.CasbinDB())
		if !errors.Is(err, nil) {
			xlog.Error("casbin New Enforcer failed",
				xlog.FieldErr(err),
				xlog.FieldComponentName("casbin"),
			)
		}
	}
	return casbinClient
}
