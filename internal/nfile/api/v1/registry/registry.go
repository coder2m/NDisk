package registry

import (
	"github.com/coder2z/g-server/xapp"
	R "github.com/coder2z/ndisk/pkg/response"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/coder2z/ndisk/internal/nfile/api/v1/handle"
	"github.com/coder2z/ndisk/internal/nfile/api/v1/middleware"
)

var (
	router *gin.Engine
	once   sync.Once
)

func init() {
	once = sync.Once{}
}

func regFileHandler(e *gin.Engine) {
	g := e.Group("/file")
	g.POST("/start", handle.Start)
	g.POST("/upload", handle.Upload)
	g.POST("/end", handle.End)
	g.GET("/download", handle.Download)
}

func Engine() *gin.Engine {
	once.Do(func() {
		if xapp.Debug() {
			gin.DisableConsoleColor()
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		router = gin.Default()
		router.NoRoute(R.HandleNotFound)
		router.Use(middleware.Auth(), middleware.GetHeader())
		regFileHandler(router)
	})
	return router
}
