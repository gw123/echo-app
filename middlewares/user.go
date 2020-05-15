package echoapp_middlewares

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func NewUserMiddlewares(skipper middleware.Skipper, usrSvr echoapp.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}
			userId, err := echoapp_util.GetCtxtUserId(c)
			if err != nil {
				echoapp_util.ExtractEntry(c).Errorf("jwsMiddleware ParseToken %s", err.Error())
				return c.JSON(http.StatusUnauthorized, "未授权")
			}
			user, err := usrSvr.GetUserById(userId)
			if err != nil {
				//这个级别比较严重，可能秘钥已经泄露(伪造用户id), 或者redis 数据库都出现问题
				echoapp_util.ExtractEntry(c).
					WithField("report", "可能秘钥已经泄露(伪造用户id), 或者redis 数据库都出现问题").
					Errorf("jwsMiddleware not found userId: %d, err: %s", userId, err.Error())
				return c.JSON(http.StatusUnauthorized, "未授权")
			}
			echoapp_util.SetCtxUser(c, user)
			return next(c)
		}
	}
}
