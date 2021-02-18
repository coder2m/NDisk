package service

import (
	"io/ioutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

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

func ConnectFile(f *model.File) error {
	switch f.FileSystem {
	case FileSystemDisk:
	case FileSystemOss:
	case FileSystemCos:
	case FileSystem7Niu:
	case FileSystemGfs:
	case FileSystemTfs:
	default:

	}
	return nil
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
