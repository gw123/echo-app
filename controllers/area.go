package controllers

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/labstack/echo"
)

type AreaController struct {
	echoapp.BaseController
}

func NewAreaController() *AreaController {
	temp := new(AreaController)
	return temp
}

func (t *AreaController) GetAreaList(ctx echo.Context) error {
	areaId := ctx.QueryParam("areaId")
	areaList, err := app.MustGetAreaService().GetAreaList(areaId)
	if err != nil {

	}
	return t.BaseController.Success(ctx, areaList)
}
