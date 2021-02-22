package registry

import (
	R "github.com/myxy99/ndisk/pkg/response"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/myxy99/ndisk/internal/nfile/api/v1/handle"
	"github.com/myxy99/ndisk/internal/nfile/api/v1/middleware"
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
}

func Engine() *gin.Engine {
	once.Do(func() {
		gin.SetMode(gin.ReleaseMode)
		router = gin.Default()
		router.NoRoute(R.HandleNotFound)
		router.Use(middleware.GetHeader())
		regFileHandler(router)
	})
	return router
}
