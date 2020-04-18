package controllers

import (
	"net/http"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/services"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type HashQrcodeContollers struct {
	echoapp.BaseController
}

func NewHashqrcontrollers() *HashQrcodeContollers {
	help := &HashQrcodeContollers{}
	return help
}
func (h *HashQrcodeContollers) GetHashQrcode(c echo.Context) error {
	code := c.QueryParam("code") //code (type string)
	// 字符串转int
	png, err := services.HashEncode(code)
	if err != nil {
		return h.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "QRcode参数错误"))
	}
	c.Response().Header().Set(echo.HeaderContentType, "image/png")
	_, err = c.Response().Write(png)

	if err != nil {
		return h.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "QRcode write"))
	}
	return nil
}

func (h *HashQrcodeContollers) Gethashdecode(c echo.Context) error {
	hashcode := c.QueryParam("hashcode")
	decodeArr, err := services.HashDecode(hashcode)
	if err != nil {
		return err
	}
	hashparam := map[string]interface{}{
		"哈希二维码：": hashcode,
		"解码序列：":  decodeArr,
	}
	return c.JSON(http.StatusOK, hashparam)

}
