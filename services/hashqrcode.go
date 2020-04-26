package services

import (
	"strconv"

	"github.com/pkg/errors"
	qrcode "github.com/skip2/go-qrcode"
	hashids "github.com/speps/go-hashids"
)

type HashCodeService struct {
}

const (
	minLength = 30
	salt      = "This is set as gh"
	size      = 255
)

func NewHashCodeService() *HashCodeService {
	return &HashCodeService{}
}
func (hash *HashCodeService) HashQrEncode(code string) ([]byte, error) {
	number, err := strconv.Atoi(code)
	if err != nil {
		return nil, err
	}
	hs := hashids.NewData()
	hs.Salt = salt
	hs.MinLength = minLength
	ha, _ := hashids.NewWithData(hs)
	e, _ := ha.Encode([]int{number})
	png, err := qrcode.Encode(e, qrcode.Medium, size)
	return png, err
}

func (hash *HashCodeService) HashDecode(code string) ([]int, error) {
	hs := hashids.NewData()
	hs.Salt = salt
	hs.MinLength = minLength
	ha, _ := hashids.NewWithData(hs)
	decodeArr, err := ha.DecodeWithError(code)
	if err != nil {
		return nil, errors.Wrap(err, "解码错误")
	}
	return decodeArr, nil
}
