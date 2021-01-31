package echoapp_middlewares

import (
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
	IsTry      bool
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
				echoapp_util.ExtractEntry(c).Info("read token from Cookie: " + token)
			}

			auth := req.Header.Get(echo.HeaderAuthorization)
			if token == "" && len(auth) == 0 && opt.IgnoreAuth {
				echoapp_util.ExtractEntry(c).Info("ignoreAuth")
				return next(c)
			}

			authScheme := "Bearer"
			l := len(authScheme)
			if token == "" {
				if !(len(auth) > l+1 && auth[:l] == authScheme) {
					echoapp_util.ExtractEntry(c).Error("未设置token")
					if opt.IsTry {
						return next(c)
					} else {
						return c.JSON(http.StatusUnauthorized, "未授权")
					}
				}
				token = auth[l+1:]
			}
			userId, payload, err := opt.Jws.ParseToken(token)
			if err != nil {
				if opt.IgnoreAuth {
					echoapp_util.ExtractEntry(c).Error("未设置token")
					return next(c)
				} else {
					echoapp_util.ExtractEntry(c).Errorf("jwsMiddleware ParseToken %s", err.Error())
					if c.Request().Header.Get(echo.HeaderContentType) == echo.MIMEApplicationJSON {
						if opt.IsTry {
							return next(c)
						} else {
							return c.JSON(http.StatusUnauthorized, "未授权2")
						}
					} else {
						if opt.IsTry {
							return next(c)
						} else {
							data := make(map[string]interface{})
							data["auth_url"] = c.Request().RequestURI
							return c.Render(http.StatusOK, "callback", data)
						}
					}
				}
			}

			echoapp_util.ExtractEntry(c).Infof("userID: %d, payload: %s", userId, payload)
			//glog.Infof("userId:%d", userId)
			echoapp_util.SetCtxUserId(c, userId)
			echoapp_util.SetCtxJwsPayload(c, payload)
			return next(c)
		}
	}
}
