package echoapp_middlewares

import (
	"context"
	"net/http"
	"time"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewWechatAuthMiddlewares(
	skipper middleware.Skipper,
	wechat echoapp.WechatService,
	userSvr echoapp.UserService,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}
			clientType := echoapp_util.GetClientTypeByUA(c.Request().UserAgent())
			comId := echoapp_util.GetCtxComId(c)
			company, err := echoapp_util.GetCtxCompany(c)
			if err != nil {
				echoapp_util.ExtractEntry(c).WithError(err).Error("getCtxCompany err")
				return next(c)
			}

			echoapp_util.ExtractEntry(c).Infof("wehchat middle Client Type:%s, openWxOfficial: %t", clientType, company.OpenWxOfficial)
			if clientType != echoapp.ClientWxOfficial || !company.OpenWxOfficial {
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

			state := c.QueryParam("state")
			echoapp_util.ExtractEntry(c).Info("params state is ", state)
			//if c.Request().URL.Path == path {
			if state == "wx_callback" {
				echoapp_util.ExtractEntry(c).Info("wx_callback")

				code := c.QueryParam("code")
				if code == "" {
					echoapp_util.ExtractEntry(c).WithError(err).Error("code为空")
					return c.HTML(http.StatusInternalServerError, "参数错误")
				}

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
					echoapp_util.SetCtxUserId(c, newUser.Id)
					echoapp_util.ExtractEntry(c).Infof("授权成功")
					c.SetCookie(&http.Cookie{
						Name:    "token",
						Value:   newUser.JwsToken,
						Expires: time.Now().Add(time.Hour * 72),
					})
				}

				return next(c)
			} else {
				//如果authtoken不存在或者校验失败， 认为用户未登录跳转到微信授权登录
				//echoapp_util.ExtractEntry(c).WithField("c.request", c.Request().URL).Infof("current page url: %s", c.Request().URL.String())
				authUrl, err := wechat.GetAuthCodeUrl(comId, c.Request().URL.Path)
				if err != nil || authUrl == "" {
					echoapp_util.ExtractEntry(c).WithError(err).Error("获取授权Url失败")
					return c.String(http.StatusInternalServerError, "系统错误请重试: not get auth url")
				}
				echoapp_util.ExtractEntry(c).Infof("jump to wxAuth authUrl %s", authUrl)
				//return c.Redirect(http.StatusFound, authUrl)
				data := make(map[string]interface{})
				data["auth_url"] = authUrl
				return c.Render(http.StatusOK, "callback", data)
			}
		}
	}
}
