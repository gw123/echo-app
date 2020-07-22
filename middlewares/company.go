package echoapp_middlewares

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"strconv"
)

func NewCompanyMiddlewares(skipper middleware.Skipper, comSvr echoapp.CompanyService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == echo.OPTIONS {
				return next(c)
			}
			if skipper(c) {
				return next(c)
			}
			comId, _ := strconv.Atoi(c.Param("com_id"))
			if comId == 0 {
				comId, _ = strconv.Atoi(c.QueryParam("com_id"))
			}
			if comId == 0 {
				echoapp_util.ExtractEntry(c).Errorf("com_id not set")
				return c.JSON(http.StatusUnauthorized, "非法请求地址")
			}
			company, err := comSvr.GetCachedCompanyById(uint(comId))
			if err != nil {
				echoapp_util.ExtractEntry(c).Errorf("com %d cache not set", comId)
				return c.JSON(http.StatusUnauthorized, "服务未初始化")
			}
			echoapp_util.SetCtxCompany(c, company)
			return next(c)
		}
	}
}
