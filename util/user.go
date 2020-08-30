package echoapp_util

import (
	"github.com/pkg/errors"
	"math/rand"
)

//获取一个加密后的userCode
func MakeUserCode(userId int64, salt string) (string, int32, error) {
	//因为这里不算写到数据库里面所以不需要考虑 索引尽可能自增加
	rand := rand.Int31n(92345678)
	userIdHash, err := EncodeInt64(userId+int64(rand), salt)
	if err != nil {
		return "", 0, errors.Wrap(err, "GetUserCode")
	}
	return userIdHash, rand, nil
}
