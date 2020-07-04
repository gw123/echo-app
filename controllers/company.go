package controllers

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
)

type CompanyController struct {
	companySvr echoapp.CompanyService
	echoapp.BaseController
}

func NewCompanyController(companySvr echoapp.CompanyService) *CompanyController {
	return &CompanyController{
		companySvr: companySvr,
	}
}

type CompanyInfoReponse struct {
	Company echoapp.CompanyBrief `json:"company"`
}

func (comCtl *CompanyController) GetCompanyInfo(ctx echo.Context) error {
	comapny, err := echoapp_util.GetCtxCompany(ctx)
	if err != nil {
		return comCtl.Fail(ctx, echoapp.CodeNotFound, "查找公司信息失败", err)
	}

	return comCtl.Success(ctx, &CompanyInfoReponse{Company: comapny.CompanyBrief})
}

func (comCtl *CompanyController) GetQuickNav(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	navs, err := comCtl.companySvr.GetQuickNav(comId)
	if err != nil {
		return comCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	return comCtl.Success(ctx, navs)
}
