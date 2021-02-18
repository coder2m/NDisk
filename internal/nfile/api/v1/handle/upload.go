package handle

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_map "github.com/myxy99/ndisk/internal/nfile/map"
	"github.com/myxy99/ndisk/internal/nfile/model"
	"github.com/myxy99/ndisk/internal/nfile/service"
	R "github.com/myxy99/ndisk/pkg/response"
	"github.com/myxy99/ndisk/pkg/utils"
)

func Start(c *gin.Context) {
	var (
		header   = GetHeader(c)
		hashType = header.HashType
		hashCode = header.HashCode
		param    = new(_map.UploadStartReq)
		file     *model.File
		err      error
	)
	//数据库查询
	if file, err = model.GetFileByHash(hashType, hashCode); err != nil && err != gorm.ErrRecordNotFound {
		R.Error(c, err.Error())
		return
	}
	//查到数据，秒传成功
	if err == nil {
		R.Ok(c, _map.UploadStartRes{
			FileId:   file.ID,
			SecTrans: true,
		})
		return
	}
	//未查到数据，初始化数据库
	if err = c.ShouldBind(param); err != nil {
		R.Error(c, err.Error())
		return
	}
	file = &model.File{
		Model: gorm.Model{},
		Name:  param.Name,
		Size:  param.Size,
		Type:  param.Type,
	}
	_ = file.SetHash(hashType, hashCode)

	service.SetBlock(file)
	service.SetFile(file)

	if err = file.Add(); err != nil {
		R.Error(c, err.Error())
		return
	}

	R.Ok(c, _map.UploadStartRes{
		FileId:    file.ID,
		SecTrans:  true,
		BlockSize: file.BlockSize,
		SliceSize: file.SliceSize(),
	})
	return
}

func Upload(c *gin.Context) {
	var (
		header    = GetHeader(c)
		file      *model.File
		sliceData []byte
		rawData   []byte
		err       error
	)
	//数据库查询
	if file, err = model.GetFileById(header.FileId); err != nil {
		R.Error(c, err.Error())
		return
	}
	if rawData, err = c.GetRawData(); err != nil {
		R.Error(c, err.Error())
		return
	}

	if sliceData, err = utils.Base64Decode(string(rawData)); err != nil {
		R.Error(c, err.Error())
		return
	}

	if len(sliceData) != int(header.Size) {
		R.Error(c, errors.New("size error"))
		return
	}

	if header.HashCode != utils.Encrypto(header.HashType, sliceData) {
		R.Error(c, errors.New("hash error"))
		return
	}
	if err = service.WriteFileSliceData(file, sliceData); err != nil {
		R.Error(c, err.Error())
		return
	}
	R.Ok(c, nil)
}

func End(c *gin.Context) {
	//header := GetHeader(c)

}
