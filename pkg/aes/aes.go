/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/6 10:18
 **/
package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func NewAes(key []byte) *cases {
	return &cases{
		key,
	}
}

type cases struct {
	key []byte
}

func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	return append(ciphertext, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	return origData[:(length - int(origData[length-1]))]
}

//AES加密
func (a cases) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	data = pKCS7Padding(data, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, a.key[:blockSize])
	encrypted := make([]byte, len(data))
	blockMode.CryptBlocks(encrypted, data)
	return encrypted, nil
}

//AES解密
func (a cases) Decrypt(encrypted []byte) ([]byte, error) {
	if len(encrypted)%aes.BlockSize != 0 {
		return nil, errors.New("decryption failed")
	}
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, a.key[:blockSize])
	origData := make([]byte, len(encrypted))
	blockMode.CryptBlocks(origData, encrypted)
	origData = pKCS7UnPadding(origData)
	return origData, nil
}
