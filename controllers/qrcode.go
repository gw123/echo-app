package controllers

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	qrcode "github.com/skip2/go-qrcode"
	"net/http"
	"strings"
)

type QrcodeController struct {
	echoapp.BaseController
}

func NewQrcodeController() *QrcodeController {
	help := &QrcodeController{
	}
	return help
}

func (h *QrcodeController) GetQrcode(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	count := strings.Count(code, "")
	if count > 64 {
		return h.Fail(ctx, echoapp.Err_Argument, "", errors.Wrap(errors.New("code 长度大于64"), "参数错误"))
	}

	png, err := qrcode.Encode(code, qrcode.Medium, 160)
	if err != nil {
		return h.Fail(ctx, echoapp.Err_Argument, "", errors.Wrap(err, "qrcode编码错误"))
	}

	ctx.Response().Header().Set(echo.HeaderContentType, "image/png")
	_, err = ctx.Response().Write(png)
	if err != nil {
		return h.Fail(ctx, echoapp.Err_Argument, "", errors.Wrap(err, "Response write"))
	}
	return nil
}

func (h *QrcodeController) Hello(ctx echo.Context) error {
	return ctx.HTML(http.StatusOK, "hello world")
}
