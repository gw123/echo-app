package echoapp_util

import (
	"encoding/base64"
	"github.com/forgoer/openssl"
)

func EntryptDesECB(data, key []byte) (string, error) {
	data, err := openssl.DesECBEncrypt(data, key, openssl.PKCS7_PADDING)
	return base64.StdEncoding.EncodeToString(data), err
}

func DecryptDESECB(data, key []byte) (string, error) {
	data, err := openssl.DesECBDecrypt(data, key, openssl.PKCS7_PADDING)
	return string(data), err
}
