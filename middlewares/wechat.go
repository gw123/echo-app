package echoapp_middlewares

import (
	"context"
	"fmt"
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/gw123/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"net/url"
)

func NewWechatAuthMiddlewares(
	skipper middleware.Skipper,
	wechat echoapp.WechatService,
	userSvr echoapp.UserService,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			echoapp_util.ExtractEntry(c).Info("UserAgent" + c.Request().UserAgent())
			if skipper(c) {
				return next(c)
			}

			comId := echoapp_util.GetCtxComId(c)
			glog.Info(c.Request().URL.Path)
			//授权回调处理
			path := fmt.Sprintf("/index-dev/%d/wxAuthCallBack", comId)
			if c.Request().URL.Path == path {
				queryValues, err := url.ParseQuery(c.Request().URL.RawQuery)
				if err != nil {
					return c.HTML(http.StatusInternalServerError, "  授权失败")
				}

				code := queryValues.Get("code")
				if code == "" {
					echoapp_util.ExtractEntry(c).WithError(err).Error("code为空")
					return c.HTML(http.StatusInternalServerError, "  授权失败")
				}
				queryState := queryValues.Get("state")
				echoapp_util.ExtractEntry(c).Info("state" + queryState)
				userInfo, err := wechat.GetUserInfo(context.Background(), comId, code)
				if err != nil {
					echoapp_util.ExtractEntry(c).WithError(err).Error("")
					return c.HTML(http.StatusInternalServerError, "  授权失败")
				}
				echoapp_util.ExtractEntry(c).Info("%+v", userInfo)
				user := &echoapp.User{
					Nickname: userInfo.Nickname,
					ComId:    comId,
					Openid:   userInfo.OpenId,
					Avatar:   userInfo.HeadImageURL,
					City:     userInfo.City,
				}
				if userInfo.Sex == 1 {
					user.Sex = "男"
				} else {
					user.Sex = "女"
				}

				echoapp_util.ExtractEntry(c).Info("自动注册用户 %+v", user)
				if err := userSvr.AutoRegisterWxUser(user); err != nil {
					echoapp_util.ExtractEntry(c).WithError(err).Error("自动注册用户失败")
					return c.HTML(http.StatusInternalServerError, "  授权失败")
				}

				return next(c)
			}

			_, err := echoapp_util.GetCtxtUserId(c)
			if err != nil {
				//如果authtoken不存在或者校验失败， 认为用户未登录跳转到微信授权登录
				authUrl, err := wechat.GetAuthCodeUrl(comId)
				if err != nil {
					echoapp_util.ExtractEntry(c).WithError(err).Error("获取授权Url失败")
					return c.String(http.StatusInternalServerError, "系统错误请重试")
				}
				return c.Redirect(http.StatusFound, authUrl)
			}
			return next(c)
		}
	}
}
