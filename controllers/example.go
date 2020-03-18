package controllers

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"net/http"
	"time"
)


type ExampleController struct {
	echoapp.BaseController
}

func NewExampleController() *ExampleController {
	help := &ExampleController{
	}
	return help
}


func (h *ExampleController) Index(ctx echo.Context) error {
	//staticParams := echoapp.Config.Asset
	renderParams := map[string]interface{}{
		"ip":  ctx.RealIP(),
		"time": time.Now().String(),
	}
	echoapp_util.ExtractEntry(ctx).Info(renderParams)
	return ctx.Render(http.StatusOK, "index", renderParams)
}
