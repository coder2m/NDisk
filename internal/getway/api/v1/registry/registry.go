package registry

import (
	"github.com/gin-gonic/gin"
	xapp "github.com/myxy99/component"
	"github.com/myxy99/ndisk/internal/getway/api/v1/middleware"
)

var router *gin.Engine

func Engine() *gin.Engine {
	if router == nil {
		if xapp.Debug() {
			gin.DisableConsoleColor()
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		router = gin.Default()
		router.Use(middleware.Cors())
	}
	return router
}

func V1() *gin.RouterGroup {
	v1 := Engine().Group("/api/v1")
	//v1.Use(middleware.CSRF())
	return v1
}
