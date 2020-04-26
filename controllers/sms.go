package controllers

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type SendMessageParams struct {
	echoapp.SendMessageOptions
}

type SmsController struct {
	echoapp.BaseController
	smsSvr echoapp.SmsService
}

func NewSmsController() *SmsController {
	smsSvr := app.MustGetSmsService()
	return &SmsController{
		smsSvr: smsSvr,
	}
}

func (sCtr *SmsController) SendMessageByToken(ctx echo.Context) error {
	params := SendMessageParams{}
	if err := ctx.Bind(&params); err != nil {
		return sCtr.Fail(ctx, echoapp.Error_ArgumentError, err.Error(), err)
	}
	echoapp_util.ExtractEntry(ctx).Infof("发送短信内容：%+v", params.SendMessageOptions)
	err := sCtr.smsSvr.SendMessage(&params.SendMessageOptions)
	if err != nil {
		return sCtr.Fail(ctx, echoapp.Error_ArgumentError, err.Error(), errors.Wrap(err, "短信发送失败,"))
	}
	return sCtr.Success(ctx, nil)
}
