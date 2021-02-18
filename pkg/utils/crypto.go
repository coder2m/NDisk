package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"strings"
)

func Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func Encrypto(htype string, data []byte) string {
	var hash hash.Hash
	switch strings.ToLower(htype) {
	case "md5":
		hash = md5.New()
	case "sha1":
		hash = sha1.New()
	case "sha256":
		hash = sha256.New()
	}
	hash.Write(data)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
