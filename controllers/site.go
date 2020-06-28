package controllers

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"strconv"
)

type SiteController struct {
	actSvr echoapp.ActivityService
	comSvr echoapp.CompanyService
	echoapp.BaseController
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

