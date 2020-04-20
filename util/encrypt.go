package echoapp_util

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"github.com/pkg/errors"
)

func EntryptDesECB(data, key []byte) (string, error) {
	if len(key) > 8 {
		key = key[:8]
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "EncryptDesECB")
	}
	bs := block.BlockSize()
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		return "", errors.New("EntryptDesECB Need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return base64.StdEncoding.EncodeToString(out), nil
}

func DecryptDESECB(d, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(string(d))
	if err != nil {
		return "", errors.Wrap(err, "DecryptDES Decode base64 error")
	}
	if len(key) > 8 {
		key = key[:8]
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "DecryptDES NewCipher error")
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return "", errors.New("DecryptDES crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = PKCS5UnPadding(out)
	return string(out), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
