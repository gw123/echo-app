package echoapp_middlewares

import (
	"github.com/gw123/glog"
	"net/http"

	"github.com/gw123/echo-app/components"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type JwsMiddlewaresOptions struct {
	Skipper middleware.Skipper
	Jws     *components.JwsHelper
	//调试时候使用直接模拟一个用户id,正式环境要把这个设置为0
	MockUserId int64
	IgnoreAuth bool
}

func NewJwsMiddlewares(opt JwsMiddlewaresOptions) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if opt.Skipper(c) {
				return next(c)
			}
			req := c.Request()

			if opt.MockUserId > 0 {
				echoapp_util.ExtractEntry(c).Infof("模拟用户:%d", opt.MockUserId)
				echoapp_util.SetCtxUserId(c, opt.MockUserId)
				echoapp_util.SetCtxJwsPayload(c, "just for test")
				return next(c)
			}

			token := c.QueryParam("token")
			tokenCookie, err := c.Cookie("token")
			if err == nil {
				token = tokenCookie.Value
				glog.Info("read token from Cookie: " + token)
			}

			auth := req.Header.Get(echo.HeaderAuthorization)
			if token == "" && len(auth) == 0 && opt.IgnoreAuth {
				return next(c)
			}

			authScheme := "Bearer"
			l := len(authScheme)
			if token == "" {
				if !(len(auth) > l+1 && auth[:l] == authScheme) {
					echoapp_util.ExtractEntry(c).Error("未设置token")
					return c.JSON(http.StatusUnauthorized, "未授权")
				}
				token = auth[l+1:]
			}
			userId, payload, err := opt.Jws.ParseToken(token)
			if err != nil {
				if opt.IgnoreAuth {
					return next(c)
				} else {
					echoapp_util.ExtractEntry(c).Errorf("jwsMiddleware ParseToken %s", err.Error())
					return c.JSON(http.StatusUnauthorized, "未授权")
				}
			}

			//glog.Infof("userId:%d", userId)
			echoapp_util.SetCtxUserId(c, userId)
			echoapp_util.SetCtxJwsPayload(c, payload)
			return next(c)
		}
	}
}
