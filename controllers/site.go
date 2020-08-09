package controllers

import (
	"fmt"
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/gw123/glog"
	"github.com/labstack/echo"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"net/http"
	"strconv"
)

type SiteController struct {
	actSvr echoapp.ActivityService
	comSvr echoapp.CompanyService
	wxSvr  echoapp.WechatService
	echoapp.BaseController
	indexCachePage []byte
}

func NewSiteController(comSvr echoapp.CompanyService, actSvr echoapp.ActivityService, svr echoapp.WechatService) *SiteController {
	return &SiteController{
		comSvr: comSvr,
		actSvr: actSvr,
		wxSvr:  svr,
	}
}

func (sCtl *SiteController) GetNotifyList(ctx echo.Context) error {
	company, err := echoapp_util.GetCtxCompany(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	notifyList, err := sCtl.actSvr.GetNotifyList(company.Id, 0, 6)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)
	}
	return sCtl.Success(ctx, notifyList)
}

func (sCtl *SiteController) GetNotifyDetail(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.QueryParam("id"))
	if id <= 0 {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", echoapp.ErrArgument)
	}
	notify, err := sCtl.actSvr.GetNotifyDetail(id)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)
	}
	return sCtl.Success(ctx, notify)
}

func (sCtl *SiteController) GetBannerList(ctx echo.Context) error {
	position := ctx.QueryParam("position")
	comId := echoapp_util.GetCtxComId(ctx)
	banner, err := sCtl.actSvr.GetBannerList(comId, position, 8)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, banner)
}

func (sCtl *SiteController) GetActivityList(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	activityList, err := sCtl.actSvr.GetActivityList(comId, lastId, limit)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, activityList)
}

func (sCtl *SiteController) GetActivityDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.QueryParam("id"))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	activity, err := sCtl.actSvr.GetActivityDetail(uint(id))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, activity)
}

func (sCtl *SiteController) GetQuickNav(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	navs, err := sCtl.comSvr.GetQuickNav(comId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	return sCtl.Success(ctx, navs)
}

func (sCtl *SiteController) Index(ctx echo.Context) error {
	//echoapp_util.ExtractEntry(ctx).Info("UserAgent" + ctx.Request().UserAgent())
	comID := echoapp_util.GetCtxComId(ctx)

	clientType := echoapp_util.GetClientTypeByUA(ctx.Request().UserAgent())
	response := make(map[string]interface{})
	response["clientType"] = clientType
	if clientType == echoapp.ClientWxOfficial {
		req := ctx.Request()
		url := fmt.Sprintf("%s://%s%s", "http", req.Host, req.URL.Path)
		jsConfig, err := sCtl.wxSvr.GetJsConfig(comID, url)
		if err != nil {
			echoapp_util.ExtractEntry(ctx).WithError(err)
		} else {
			response["wxCfg"] = jsConfig
		}
	}
	return ctx.Render(http.StatusOK, "index", response)
}

func (sCtl *SiteController) WxAuthCallBack(ctx echo.Context) error {
	comID := echoapp_util.GetCtxComId(ctx)
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return ctx.HTML(502, "授权失败")
	}
	data := make(map[string]interface{})
	data["userToken"] = user.JwsToken
	data["nickname"] = user.Nickname
	data["avatar"] = user.Avatar
	data["sex"] = user.Sex
	data["roles"] = user.Roles
	data["id"] = user.Id
	clientType := echoapp_util.GetClientTypeByUA(ctx.Request().UserAgent())
	response := make(map[string]interface{})
	response["clientType"] = clientType
	response["user"] = data
	//response["assetHost"] = "http://192.168.187.1:8889"
	response["assetHost"] = "http://m.xytschool.com/dev/public"
	if clientType == echoapp.ClientWxOfficial {
		req := ctx.Request()
		url := ""
		//if req.URL.Scheme == "" {
		//
		//}
		//todo
		if req.URL.RawQuery != "" {
			url = fmt.Sprintf("%s://%s%s?%s", "http", req.Host, req.URL.Path, req.URL.RawQuery)
		} else {
			url = fmt.Sprintf("%s://%s%s", "http", req.Host, req.URL.Path)
		}
		jsConfig, err := sCtl.wxSvr.GetJsConfig(comID, url)
		if err != nil {
			echoapp_util.ExtractEntry(ctx).WithError(err)
		} else {
			response["wxCfg"] = jsConfig
		}
	}
	return ctx.Render(http.StatusOK, "index", response)
}

func (sCtl *SiteController) WxMessage(ctx echo.Context) error {
	comID := echoapp_util.GetCtxComId(ctx)
	server, err := sCtl.wxSvr.GetOfficialServer(ctx, comID)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeInnerError, err.Error(), err)
	}
	server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
		//TODO
		//回复消息：演示回复用户发送的消息
		glog.Info(msg.Content)
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})
	if err = server.Serve(); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeInnerError, err.Error(), err)
	}
	server.Send()
	return nil
}
