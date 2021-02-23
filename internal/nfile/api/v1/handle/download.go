package handle

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/myxy99/ndisk/internal/nfile/model"
	"github.com/myxy99/ndisk/internal/nfile/service"
	R "github.com/myxy99/ndisk/pkg/response"
	"github.com/myxy99/ndisk/pkg/utils"
)

func Download(c *gin.Context) {
	var (
		header    = GetHeader(c)
		file      *model.File
		fileSilce model.FileSlice
		data      []byte
		err       error
	)
	//数据库查询
	if file, err = model.GetFileById(header.FileId); err != nil {
		R.Error(c, err.Error())
		return
	}

	if header.SliceIndex == 0 {
		R.Ok(c, file)
		return
	}

	if fileSilce, err = file.GetSlice(uint(header.SliceIndex)); err != nil {
		R.Error(c, err.Error())
		return
	}

	if data, err = service.FileGet(file.FileSystem, file.TmpFilePath(header.SliceIndex)); err != nil {
		R.Error(c, err.Error())
		return
	}

	c.Header("file_id", fmt.Sprintf("%d", fileSilce.FileId))
	c.Header("slice_index", fmt.Sprintf("%d", fileSilce.Index))
	c.Header("size", fmt.Sprintf("%d", fileSilce.Size))
	c.Header("hash_type", fileSilce.HashType)
	c.Header("hash_code", fileSilce.HashCode)

	c.String(http.StatusOK, "slice-download", utils.Base64DEncode(data))

}
