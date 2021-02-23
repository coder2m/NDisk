package handle

import (
	"fmt"
	"io"

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

	if file.Status == 0 {
		R.Error(c, "file Upload not OK")
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

	c.Header("File-Id", fmt.Sprintf("%d", fileSilce.FileId))
	c.Header("Slice-Index", fmt.Sprintf("%d", fileSilce.Index))
	c.Header("Size", fmt.Sprintf("%d", fileSilce.Size))
	c.Header("Hash-Type", fileSilce.HashType)
	c.Header("Hash-Code", fileSilce.HashCode)
	c.Header("Content-Type", "slice-download")

	io.WriteString(c.Writer, utils.Base64DEncode(data))
}
