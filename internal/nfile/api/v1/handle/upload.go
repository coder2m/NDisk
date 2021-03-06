package handle

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_map "github.com/coder2z/ndisk/internal/nfile/map"
	"github.com/coder2z/ndisk/internal/nfile/model"
	"github.com/coder2z/ndisk/internal/nfile/service"
	R "github.com/coder2z/ndisk/pkg/response"
	"github.com/coder2z/ndisk/pkg/utils"
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
	uid := c.GetInt64("Uid")
	file = &model.File{
		Model:   gorm.Model{},
		Name:    param.Name,
		Size:    param.Size,
		Type:    param.Type,
		Creator: uint64(uid),
	}
	_ = file.SetHash(hashType, hashCode)

	service.SetBlock(file)
	service.SetFile(file)
	file.SliceSize()

	if err = file.Add(); err != nil {
		R.Error(c, err.Error())
		return
	}

	R.Ok(c, _map.UploadStartRes{
		FileId:    file.ID,
		SecTrans:  false,
		BlockSize: file.BlockSize,
		SliceSize: file.SliceCount,
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

	if file.CheckSlice(uint(header.SliceIndex)) > 0 {
		R.Error(c, "slice exists")
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

	if header.HashCode != utils.Encrypt(header.HashType, sliceData) {
		R.Error(c, errors.New("hash error"))
		return
	}

	slice := file.NewSlice(uint(header.SliceIndex), header.HashType, header.HashCode, uint64(header.Size))
	if err = slice.Add(); err != nil {
		R.Error(c, err.Error())
		return
	}

	if err = service.WriteFileSliceData(file, sliceData, header.SliceIndex); err != nil {
		R.Error(c, err.Error())
		return
	}
	R.Ok(c, nil)
}

func End(c *gin.Context) {
	var (
		header = GetHeader(c)
		file   *model.File
		err    error
	)
	//数据库查询
	if file, err = model.GetFileById(header.FileId); err != nil {
		R.Error(c, err.Error())
		return
	}
	if err = service.MergeFile(file); err != nil {
		R.Error(c, err.Error())
		return
	}

	R.Ok(c, nil)
}
