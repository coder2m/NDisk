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

func V1() *gin.RouterGroup {
	v1 := Engine().Group("/api/v1")
	//v1.Use(middleware.CSRF())
	return v1
}
