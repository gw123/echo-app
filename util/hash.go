package echoapp_util

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
)

func Sha256(src string) string {
	return encode(sha256.New(), []byte(src))
}

func Md5(src string) string {
	return encode(md5.New(), []byte(src))
}

func encode(h hash.Hash, src []byte) string {
	return fmt.Sprintf("%x", md5.Sum(src))
}
