package service

import (
	"bytes"
	"errors"

	"github.com/coder2z/g-server/xinvoker/oss/standard"

	xclient "github.com/coder2z/ndisk/internal/nfile/client"
	"github.com/coder2z/ndisk/internal/nfile/model"
	"github.com/coder2z/ndisk/pkg/utils"
)

const (
	FileSystemDisk = iota
	FileSystemOss
	FileSystemCos
	FileSystem7Niu
	FileSystemGfs
	FileSystemTfs
)

const (
	CONV = 1024
	B    = 1
	KB   = CONV * B
	MB   = CONV * KB
	GB   = CONV * MB
	TB   = CONV * GB
)

func WriteFileSliceData(f *model.File, data []byte, idx int) error {
	return FileSave(f.FileSystem, f.TmpFilePath(idx), data)
}

func MergeFile(f *model.File) error {
	data := make([]byte, 0)
	for i := 0; i < int(f.SliceCount); i++ {
		tmp, err := FileGet(f.FileSystem, f.TmpFilePath(i))
		if err != nil {
			return err
		}
		data = append(data, tmp...)
	}

	if len(data) != int(f.Size) {
		return errors.New("size validate error")
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
	if err := f.SetStatus(1); err != nil {
		return err
	}
	return FileSave(f.FileSystem, f.TmpMergeFilePath(), data)
}

func FileClient(sys uint8) standard.Oss {
	var client standard.Oss
	switch sys {
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
		panic("file system not find")
	}
	return client
}

func FileSave(sys uint8, path string, data []byte) error {
	return FileClient(sys).PutObject(path, bytes.NewReader(data))
}

func FileGet(sys uint8, path string) (data []byte, err error) {
	return FileClient(sys).GetObject(path)
}

func CurrBlock() uint64 {
	return MB
}
func CurrentFileSystem() uint8 {
	return FileSystemOss
}
func SetBlock(f *model.File) {
	f.BlockSize = CurrBlock()
}

func SetFile(f *model.File) {
	f.FileSystem = CurrentFileSystem()
}
