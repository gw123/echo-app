package controllers

import (
	"net/http"
	"strconv"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/gw123/glog"
	"github.com/labstack/echo"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type SiteController struct {
	actSvr    echoapp.ActivityService
	bannerSvr echoapp.SiteService
	comSvr    echoapp.CompanyService
	wxSvr     echoapp.WechatService
	echoapp.BaseController
	assetOpts      echoapp.Asset
	indexCachePage []byte
	videoSvr       echoapp.VideoService
}

func NewSiteController(comSvr echoapp.CompanyService,
	actSvr echoapp.ActivityService,
	bannerSvr echoapp.SiteService,
	svr echoapp.WechatService,
	video echoapp.VideoService,
	asset echoapp.Asset,
) *SiteController {
	return &SiteController{
		comSvr:    comSvr,
		actSvr:    actSvr,
		bannerSvr: bannerSvr,
		wxSvr:     svr,
		videoSvr:  video,
		assetOpts: asset,
	}
}

func (sCtl *SiteController) GetNotifyList(ctx echo.Context) error {
	company, err := echoapp_util.GetCtxCompany(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	notifyList, err := sCtl.bannerSvr.GetNotifyList(company.Id, 0, 6)
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
	notify, err := sCtl.bannerSvr.GetNotifyDetail(id)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)
	}
	return sCtl.Success(ctx, notify)
}

func (sCtl *SiteController) GetIndexPageBanners(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	pageBanners, err := sCtl.bannerSvr.GetIndexPageBanners(comId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, pageBanners)
}

func (sCtl *SiteController) GetBannerList(ctx echo.Context) error {
	position := ctx.QueryParam("position")
	page := ctx.QueryParam("page")
	comId := echoapp_util.GetCtxComId(ctx)
	banner, err := sCtl.bannerSvr.GetBannerList(comId, page, position, 8)
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
	clientType := echoapp_util.GetClientTypeByUA(ctx.Request().UserAgent())
	response := make(map[string]interface{})
	response["clientType"] = clientType
	response["publicHost"] = echoapp_util.GetOptimalPublicHost(ctx, sCtl.assetOpts)
	if clientType == echoapp.ClientWxOfficial {
		user, err := echoapp_util.GetCtxtUser(ctx)
		if err == nil {
			data := make(map[string]interface{})
			data["userToken"] = user.JwsToken
			data["nickname"] = user.Nickname
			data["avatar"] = user.Avatar
			data["sex"] = user.Sex
			data["roles"] = user.Roles
			data["id"] = user.Id
			data["vip_level"] = user.VipLevel
			response["user"] = data
		}
	}
	return ctx.Render(http.StatusOK, "index", response)
}

//todo 这里有302无限循环的风险
func (sCtl *SiteController) WxAuthCallBack(ctx echo.Context) error {
	//user, err := echoapp_util.GetCtxtUser(ctx)
	//if err != nil {
	//	return ctx.HTML(502, "授权失败")
	//}
	//ctx.SetCookie(&http.Cookie{
	//	Name:    "token",
	//	Value:   user.JwsToken,
	//	Expires: time.Now().Add(time.Hour * 24 * 30),
	//})
	return ctx.Render(http.StatusOK, "callback", nil)
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

func (sCtl *SiteController) GetWxConfig(ctx echo.Context) error {
	comID := echoapp_util.GetCtxComId(ctx)
	url := ctx.QueryParam("url")
	glog.DefaultLogger().Info("GetWxConfig: " + url)
	jsConfig, err := sCtl.wxSvr.GetJsConfig(comID, url)
	if err != nil {
		echoapp_util.ExtractEntry(ctx).WithError(err).Error("获取JSconfig失败")
		return sCtl.AppErr(ctx, echoapp.AppErrArgument)
	}
	return sCtl.Success(ctx, jsConfig)
}

func (sCtl *SiteController) GetVideoList(ctx echo.Context) error {
	comID := echoapp_util.GetCtxComId(ctx)
	lastID, limit := echoapp_util.GetCtxListParams(ctx)
	videoList, err := sCtl.videoSvr.GetVideoList(comID, lastID, limit)
	if err != nil {
		echoapp_util.ExtractEntry(ctx).WithError(err).Error("获取Video失败")
		return sCtl.AppErr(ctx, echoapp.AppErrArgument)
	}
	return sCtl.Success(ctx, videoList)
}

func (sCtl *SiteController) GetVideoDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	video, err := sCtl.videoSvr.GetVideoDetail(uint(id))
	if err != nil {
		echoapp_util.ExtractEntry(ctx).WithError(err).Error("获取Video详情失败")
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return ctx.Render(http.StatusOK, "video", video)
}
