package controllers

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/labstack/echo"
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
	err := sCtr.smsSvr.SendMessage(params.SendMessageOptions)
	if err != nil {
		return sCtr.Fail(ctx, echoapp.Error_ArgumentError, err.Error(), err)
	}
	return sCtr.Success(ctx, nil)
}
