package service

import (
	"errors"
	"io/ioutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/myxy99/component/xinvoker/oss/standard"

	xclient "github.com/myxy99/ndisk/internal/nfile/client"
	"github.com/myxy99/ndisk/internal/nfile/model"
	"github.com/myxy99/ndisk/pkg/utils"
)

const (
	FileSystemDisk = iota
	FileSystemOss
	FileSystemCos
	FileSystem7Niu
	FileSystemGfs
	FileSystemTfs
)

func WriteFileSliceData(f *model.File, data []byte, idx int) error {
	return ioutil.WriteFile(f.TmpFilePath(idx), data, oss.FilePermMode)
}

func OnlyMergeFile(f *model.File) error {
	data := make([]byte, 0)
	for i := 0; i < int(f.SliceCount); i++ {
		tmp, err := ioutil.ReadFile(f.TmpFilePath(i))
		if err != nil {
			return err
		}
		data = append(data, tmp...)
	}

	if len(data) != int(f.Size) {
		return errors.New("size valide error")
	}
	var (
		hashType string
		hashCode string
	)
	switch {
	case len(f.Md5) > 0:
		hashType = "md5"
		hashCode = f.Md5
	case len(f.Sha1) > 0:
		hashType = "sha1"
		hashCode = f.Sha1
	case len(f.Sha256) > 0:
		hashType = "sha256"
		hashCode = f.Sha256
	}
	if utils.Encrypt(hashType, data) != hashCode {
		return errors.New(hashType + " validate error")
	}
	return ioutil.WriteFile(f.TmpMergeFilePath(), data, oss.FilePermMode)
}

func MergeFile(f *model.File) error {
	if err := OnlyMergeFile(f); err != nil {
		return err
	}
	var client standard.Oss
	switch f.FileSystem {
	case FileSystemDisk:
		client = xclient.Disk()
	case FileSystemOss:
		client = xclient.Oss()
	case FileSystemCos:
		client = xclient.Cos()
	case FileSystem7Niu:
		client = xclient.SevenNiu()
	case FileSystemGfs:
		client = xclient.Gfs()
	case FileSystemTfs:
		client = xclient.Tfs()
	default:
		return errors.New("file system not find")
	}
	return client.PutObjectFromFile(f.FileRealPath, f.TmpMergeFilePath())
}

func CurrBlock() uint64 {
	return 10
}
func CurrentFileSystem() uint8 {
	return 0
}
func CurrentFileSystemBasePath() string {
	return ""
}

func SetBlock(f *model.File) {
	f.BlockSize = CurrBlock()
}

func SetFile(f *model.File) {
	f.FileSystem = CurrentFileSystem()
	f.FileRealPath = CurrentFileSystemBasePath()
}
