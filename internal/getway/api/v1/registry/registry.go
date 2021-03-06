package registry

import (
	R "github.com/coder2z/ndisk/pkg/response"
	"time"

	xapp "github.com/coder2z/component"
	"github.com/coder2z/ndisk/internal/getway/api/v1/middleware"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
	v1     *gin.RouterGroup
)

func Engine() *gin.Engine {
	if router == nil {
		if xapp.Debug() {
			gin.DisableConsoleColor()
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		router = gin.New()
		router.NoRoute(R.HandleNotFound)
		router.Use(
			middleware.RecoverMiddleware(20*time.Second),
			middleware.XMonitor(),
			middleware.XTrace(),
		)
	}
	return router
}

func V1() *gin.RouterGroup {
	if v1 == nil {
		r := Engine()
		v1 = r.Group("/api/v1")
		v1.Use(
			middleware.Cors(),
			//middleware.CSRF(),
		)
	}
	return v1
}
