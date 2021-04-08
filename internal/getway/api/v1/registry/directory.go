package registry

import (
	"github.com/coder2z/ndisk/internal/getway/api/v1/handle/directory"
	"github.com/coder2z/ndisk/internal/getway/api/v1/middleware"
)

func init() {
	dir := V1().Group("/directory")
	dir.Use(middleware.Auth())
	dir.GET("/:id", directory.List)
	dir.POST("/", directory.Add)
	dir.DELETE("/:id", directory.Del)
	dir.PUT("/:id", directory.Update)
}
