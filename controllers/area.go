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

func (c *AreaController) GetAreaMap(ctx echo.Context) error {
	areaId := ctx.QueryParam("area_id")
	areaMap, err := app.MustGetAreaService().GetAreaMap(areaId)
	if err != nil {
		return c.BaseController.Fail(ctx, echoapp.Err_Argument, "", err)
	}
	return c.BaseController.Success(ctx, areaMap)
}

func (c *AreaController) GetAreaArray(ctx echo.Context) error {
	areaId := ctx.QueryParam("area_id")
	areaArray, err := app.MustGetAreaService().GetAreaArray(areaId)
	if err != nil {
		return c.BaseController.Fail(ctx, echoapp.Err_Argument, "", err)
	}
	return c.BaseController.Success(ctx, areaArray)
}
