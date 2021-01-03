package registry

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func Engine() *gin.Engine {
	if router == nil {
		gin.SetMode(gin.ReleaseMode)
		router = gin.Default()
	}
	return router
}
