package controllers

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"strconv"
)

type GoodsController struct {
	goodsSvr echoapp.GoodsService
	echoapp.BaseController
}

func NewGoodsController(usrSvr echoapp.GoodsService) *GoodsController {
	return &GoodsController{
		goodsSvr: usrSvr,
	}
}

func (sCtl *GoodsController) GetIndexBanners(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	banners, err := sCtl.goodsSvr.GetIndexBanner(comId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	return sCtl.Success(ctx, banners)
}

func (sCtl *GoodsController) GetGoodsList(ctx echo.Context) error {
	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	comId := echoapp_util.GetCtxComId(ctx)
	goods, err := sCtl.goodsSvr.GetGoodsList(comId, lastId, limit)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	return sCtl.Success(ctx, goods)
}

func (sCtl *GoodsController) GetRecommendGoods(ctx echo.Context) error {
	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	comId := echoapp_util.GetCtxComId(ctx)
	goods, err := sCtl.goodsSvr.GetRecommendGoodsList(comId, lastId, limit)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	return sCtl.Success(ctx, goods)
}

func (sCtl *GoodsController) GetGoodsInfo(ctx echo.Context) error {
	goodsId, err := strconv.Atoi(ctx.QueryParam("goodsId"))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	goods, err := sCtl.goodsSvr.GetGoodsById(goodsId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	return sCtl.Success(ctx, goods)
}


