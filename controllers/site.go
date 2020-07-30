package controllers

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/gw123/glog"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type SiteController struct {
	actSvr echoapp.ActivityService
	comSvr echoapp.CompanyService
	echoapp.BaseController
	indexCachePage []byte
}

func NewSiteController(comSvr echoapp.CompanyService, actSvr echoapp.ActivityService) *SiteController {
	return &SiteController{
		comSvr: comSvr,
		actSvr: actSvr,
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
	echoapp_util.ExtractEntry(ctx).Info("UserAgent" + ctx.Request().UserAgent())
	//if len(sCtl.indexCachePage) == 0 {
	//	indexFilePath := echoapp.ConfigOpts.Asset.PublicRoot + "/m/index.html"
	//	var err error
	//	sCtl.indexCachePage, err = ioutil.ReadFile(indexFilePath)
	//	if err != nil {
	//		return ctx.HTML(502, "文件不存在")
	//	}
	//}
	return ctx.Render(http.StatusOK, "index", nil)
}

func (sCtl *SiteController) WxAuthCallBack(ctx echo.Context) error {
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return ctx.HTML(502, "授权失败")
	}
	data := make(map[string]interface{})
	glog.Infof("WxAUthCallback UserInfo :+v", user)
	data["userToken"] = user.JwsToken
	data["nickname"] = user.Nickname
	data["avatar"] = user.Avatar
	data["sex"] = user.Sex
	data["roles"] = user.Roles
	data["id"] = user.Id

	return ctx.Render(http.StatusOK, "index", data)
}
