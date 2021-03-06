package model

import (
	"fmt"
	"strings"

	"github.com/coder2z/component/xcfg"
	"github.com/coder2z/g-saber/xlog"
	"gorm.io/gorm"
)

type (
	File struct {
		gorm.Model
		Name       string
		Size       uint64
		Type       string
		Md5        string
		Sha1       string
		Sha256     string
		FileSystem uint8
		BlockSize  uint64
		Creator    uint64
		SliceCount uint64
		Status     uint8
	}

	FileSlice struct {
		gorm.Model
		FileId   uint
		HashType string
		HashCode string
		Size     uint64
		Index    uint
	}
)

func (f *File) TableName() string {
	return "db_common_file"
}

func (f *FileSlice) TableName() string {
	return "db_common_file_slice"
}

func (f *File) SetHash(ht, hc string) (err error) {
	switch strings.ToLower(ht) {
	case "md5":
		f.Md5 = hc
	case "sha1":
		f.Sha1 = hc
	case "sha256":
		f.Sha256 = hc
	default:
		err = fmt.Errorf("can not support hashType %s", ht)
	}
	return
}

func (f *File) Add() (err error) {
	err = MainDB().Create(f).Error
	return
}

func (f *File) SliceSize() {
	f.SliceCount = (f.Size + f.BlockSize - 1) / f.BlockSize
}

func (f *File) TmpFilePath(idx int) string {
	return fmt.Sprintf("%s/%d/%d/%d", xcfg.GetString("tmp_file_path"), f.FileSystem, f.ID, idx)
}
func (f *File) TmpMergeFilePath() string {
	return fmt.Sprintf("%s/%d/%d/final_file", xcfg.GetString("tmp_file_path"), f.FileSystem, f.ID)
}
func (f *File) NewSlice(idx uint, hashType, hashCode string, size uint64) *FileSlice {
	return &FileSlice{
		FileId:   f.ID,
		HashType: hashType,
		HashCode: hashCode,
		Size:     size,
		Index:    idx,
	}
}

func (f *File) CheckSlice(idx uint) (count int64) {
	fs := &FileSlice{}
	if err := MainDB().Table(fs.TableName()).Where("file_id=? and index=?", f.ID, idx).Count(&count).Error; err != nil {
		xlog.Errorw("checkout slice ", "error", err.Error())
	}
	return
}

func (f *File) GetSlice(idx uint) (slice FileSlice, err error) {
	if err = MainDB().Table(slice.TableName()).Where("file_id=? and index=?", f.ID, idx).Take(&slice).Error; err != nil {
		xlog.Errorw("get slice ", "error", err.Error())
	}
	return
}

func (f *FileSlice) Add() (err error) {
	err = MainDB().Create(f).Error
	return
}

func (f *File) SetStatus(status uint8) (err error) {
	err = MainDB().Table(f.TableName()).Where("id=?", f.ID).Update("status", status).Error
	return
}
