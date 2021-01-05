package registry

import (
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func Engine() *gin.Engine {
	if Router == nil {
		gin.SetMode(gin.ReleaseMode)
		Router = gin.Default()
	}
	return Router
}

