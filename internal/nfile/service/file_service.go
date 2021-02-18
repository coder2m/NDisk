package service

import "github.com/myxy99/ndisk/internal/nfile/model"

func WriteFileSliceData(f *model.File, data []byte) error {
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
