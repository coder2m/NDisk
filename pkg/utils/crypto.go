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

func Encrypt(htype string, data []byte) string {
	var h hash.Hash
	switch strings.ToLower(htype) {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	default:
		panic("type unknown")
	}
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
