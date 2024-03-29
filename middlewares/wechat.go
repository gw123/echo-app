package echoapp_middlewares

import (
	"context"
	"fmt"
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"net/url"
	"strings"
)

func NewWechatAuthMiddlewares(
	skipper middleware.Skipper,
	wechat echoapp.WechatService,
	userSvr echoapp.UserService,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//echoapp_util.ExtractEntry(c).Info("UserAgent" + c.Request().UserAgent())
			if skipper(c) {
				return next(c)
			}
			clientType := echoapp_util.GetClientTypeByUA(c.Request().UserAgent())

			//echoapp_util.ExtractEntry(c).Info("Client Type:" + clientType)
			if clientType != echoapp.ClientWxOfficial {
				return next(c)
			}

			userId, err := echoapp_util.GetCtxtUserId(c)
			//echoapp_util.ExtractEntry(c).Infof("User err: %+v, user: %+v", err, userId)
			if err == nil && userId != 0 {
				// 如果用户已经登陆的逻辑
				// 补偿机制当没有调用user中间件时候 主动获取用户
				user, err := echoapp_util.GetCtxtUser(c)
				if err == nil {
					return next(c)
				} else {
					if user, err = userSvr.GetUserById(userId); err == nil {
						echoapp_util.SetCtxUser(c, user)
						return next(c)
					} else {
						echoapp_util.ExtractEntry(c).WithError(err).Errorf("获取用户信息失败UserId:%d", userId)
					}
				}
			}

			comId := echoapp_util.GetCtxComId(c)
			//授权回调处理
			var path string
			if strings.HasPrefix(c.Request().URL.Path, "/index-dev") {
				path = fmt.Sprintf("/index-dev/%d/wxAuthCallBack", comId)
			} else {
				path = fmt.Sprintf("/index/%d/wxAuthCallBack", comId)
			}

			if c.Request().URL.Path == path {
				queryValues, err := url.ParseQuery(c.Request().URL.RawQuery)
				if err != nil {
					return c.HTML(http.StatusInternalServerError, "url解析失败")
				}

				code := queryValues.Get("code")
				if code == "" {
					echoapp_util.ExtractEntry(c).WithError(err).Error("code为空")
					return c.HTML(http.StatusInternalServerError, "参数错误")
				}
				//queryState := queryValues.Get("state")
				//echoapp_util.ExtractEntry(c).Info("state" + queryState)
				userInfo, err := wechat.GetUserInfo(context.Background(), comId, code)
				if err != nil {
					echoapp_util.ExtractEntry(c).WithError(err).Error("微信授权失败")
					return c.HTML(http.StatusInternalServerError, "授权失败")
				}

				//echoapp_util.ExtractEntry(c).Info("wechatUser: %+v", userInfo)
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

				if newUser, err := userSvr.AutoRegisterWxUser(user); err != nil {
					echoapp_util.ExtractEntry(c).WithError(err).Error("自动注册用户失败")
					return c.HTML(http.StatusInternalServerError, "注册新用户失败")
				} else {
					echoapp_util.SetCtxUser(c, newUser)
					echoapp_util.ExtractEntry(c).Infof("自动登录用户-> %+v", newUser)
					//token := fmt.Sprintf("?token=%s", newUser.JwsToken)
					//return c.Redirect(http.StatusMovedPermanently, strings.Replace(path, "/wxAuthCallBack", token, -1))
				}
				return next(c)
			} else {
				//如果authtoken不存在或者校验失败， 认为用户未登录跳转到微信授权登录
				authUrl, err := wechat.GetAuthCodeUrl(comId)
				if err != nil || authUrl == "" {
					echoapp_util.ExtractEntry(c).WithError(err).Error("获取授权Url失败")
					return c.String(http.StatusInternalServerError, "系统错误请重试: not get auth url")
				}
				return c.Redirect(http.StatusFound, authUrl)
			}
		}
	}
}
