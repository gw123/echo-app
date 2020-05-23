package controllers

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/services"
	"github.com/labstack/echo"
)

type PrintController struct {
	printSvr echoapp.PrintService
	echoapp.BaseController
}

func NewPrintController() *PrintController {
	return &PrintController{
		printSvr: &services.PrintService{},
	}
}

func (pCtl *PrintController) PrintTicket(ctx echo.Context) error {
	type Params struct {
		Code        string `json:"code"`
		Username    string `json:"username"`
		ComId       int    `json:"com_id"`
		PrinterName string `json:"username"`
	}

	params := &Params{}
	if err := ctx.Bind(params); err != nil {
		return pCtl.Fail(ctx, echoapp.Err_Argument, err.Error(), err)
	}

	if params.PrinterName == "" {
		params.PrinterName = "ticket"
	}

	if err := pCtl.printSvr.PrintByPrinterName(params.ComId, params.PrinterName, []byte(params.Code)); err != nil {
		return pCtl.Fail(ctx, echoapp.Err_InnerError, err.Error(), err)
	}
	return pCtl.Success(ctx, nil)
}
