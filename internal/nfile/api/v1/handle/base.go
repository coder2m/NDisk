package handle

import (
	"strings"

	"github.com/gin-gonic/gin"

	_map "github.com/myxy99/ndisk/internal/nfile/map"
	"github.com/myxy99/ndisk/internal/nfile/model"
)

func GetHeader(c *gin.Context) *_map.Header {
	return c.MustGet("file_header").(*_map.Header)
}

func SetHash(ht, hc string, f *model.File) {
	switch strings.ToLower(ht) {
	case "md5":
		f.Md5 = hc
	case "sha1":
		f.Sha1 = hc
	case "sha256":
		f.Sha256 = hc
	}
}
