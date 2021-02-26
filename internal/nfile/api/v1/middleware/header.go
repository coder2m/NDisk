package middleware

import (
	"github.com/gin-gonic/gin"

	_map "github.com/coder2m/ndisk/internal/nfile/map"
	R "github.com/coder2m/ndisk/pkg/response"
	"github.com/coder2m/ndisk/pkg/utils"
)

func GetHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			h   = _map.NewHeader()
			err error
		)
		if err = utils.BindHttpHeader(c, h); err == nil {
			err = h.Validate()
		}
		if err != nil {
			R.HandleBadRequest(c, err.Error())
			c.Abort()
		}

		c.Set("file_header", h)
	}
}
