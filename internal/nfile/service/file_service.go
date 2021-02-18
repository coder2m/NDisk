package service

import (
	"errors"
	"io/ioutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/myxy99/component/xinvoker/oss/standard"

	xclient "github.com/myxy99/ndisk/internal/nfile/client"
	"github.com/myxy99/ndisk/internal/nfile/model"
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

	return nil
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
