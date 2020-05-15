package echoapp_middlewares

import (
	"github.com/gw123/echo-app/components"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

const HeaderToken = "X-Token"

func NewJwsMiddlewares(skipper middleware.Skipper, jws *components.JwsHelper) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}
			req := c.Request()
			token := req.Header.Get(HeaderToken)
			userId, payload, err := jws.ParseToken(token)
			if err != nil {
				echoapp_util.ExtractEntry(c).Errorf("jwsMiddleware ParseToken %s", err.Error())
				return c.JSON(http.StatusUnauthorized, "未授权")
			}
			echoapp_util.SetCtxUserId(c, userId)
			echoapp_util.SetCtxJwsPayload(c, payload)
			return next(c)
		}
	}
}
