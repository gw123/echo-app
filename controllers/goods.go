package controllers

import (
	"strconv"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
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

//func (sCtl *GoodsController) GetIndexBanners(ctx echo.Context) error {
//	comId := echoapp_util.GetCtxComId(ctx)
//	banners, err := sCtl.goodsSvr.GetIndexBanner(comId)
//	if err != nil {
//		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
//	}
//	return sCtl.Success(ctx, banners)
//}

func (sCtl *GoodsController) GetGoodsList(ctx echo.Context) error {
	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	comId := echoapp_util.GetCtxComId(ctx)
	keyword := ctx.QueryParam("keyword")
	var goods []*echoapp.GoodsBrief
	var err error
	if keyword == "" {
		goods, err = sCtl.goodsSvr.GetGoodsList(comId, lastId, limit)
		if err != nil {
			return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
		}
	} else {
		goods, err = sCtl.goodsSvr.GetGoodsListByKeyword(comId, keyword, lastId, limit)
		if err != nil {
			return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
		}
	}

	return sCtl.Success(ctx, goods)
}

func (sCtl *GoodsController) GetTagGoodsList(ctx echo.Context) error {
	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	comId := echoapp_util.GetCtxComId(ctx)
	tagID, err := strconv.Atoi(ctx.QueryParam("tag_id"))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	goodsList, err := sCtl.goodsSvr.GetTagGoodsList(comId, tagID, lastId, limit)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	return sCtl.Success(ctx, goodsList)
}

func (sCtl *GoodsController) GetGoodsTags(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	goodsList, err := sCtl.goodsSvr.GetGoodsTags(comId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	return sCtl.Success(ctx, goodsList)
}

func (sCtl *GoodsController) GetRecommendGoodsList(ctx echo.Context) error {
	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	comId := echoapp_util.GetCtxComId(ctx)
	goods, err := sCtl.goodsSvr.GetRecommendGoodsList(comId, lastId, limit)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	return sCtl.Success(ctx, goods)
}

func (sCtl *GoodsController) GetGoodsInfo(ctx echo.Context) error {
	goodsId, err := strconv.Atoi(ctx.QueryParam("id"))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	goods, err := sCtl.goodsSvr.GetGoodsById(uint(goodsId))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	if err := sCtl.goodsSvr.AddGoodsPv(goodsId); err != nil {
		echoapp_util.ExtractEntry(ctx).Error(err)
	}
	goods.Pv += 1
	return sCtl.Success(ctx, goods)
}

// 用户办理会员时候使用的商品
func (sCtl *GoodsController) GetVipDesc(ctx echo.Context) error {
	goods, err := sCtl.goodsSvr.GetVipDesc()
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现商品", err)
	}
	return sCtl.Success(ctx, goods)
}

func (sCtl *GoodsController) GetCartGoodsList(ctx echo.Context) error {
	userID, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "授权失败", err)
	}
	comID := echoapp_util.GetCtxComId(ctx)

	cart, err := sCtl.goodsSvr.GetCartGoodsList(comID, uint(userID))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "获取购物车失败", err)
	}
	return sCtl.Success(ctx, cart.Content)
}

func (sCtl *GoodsController) AddCartGoods(ctx echo.Context) error {
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "授权失败", err)
	}
	comID := echoapp_util.GetCtxComId(ctx)
	var goosdItem *echoapp.CartGoodsItem
	if err := ctx.Bind(&goosdItem); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}
	//强制新加时数量为1
	goosdItem.Num = 1
	if err := sCtl.goodsSvr.AddCartGoods(comID, uint(userId), goosdItem); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "保存失败", err)
	}
	return sCtl.Success(ctx, nil)
}

func (sCtl *GoodsController) DelCartGoods(ctx echo.Context) error {
	userID, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "授权失败", err)
	}
	comID := echoapp_util.GetCtxComId(ctx)
	var goosdItem *echoapp.CartGoodsItem
	if err := ctx.Bind(&goosdItem); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	if err := sCtl.goodsSvr.DelCartGoods(comID, uint(userID), goosdItem.GoodsId, goosdItem.SkuID); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "保存失败", err)
	}
	return sCtl.Success(ctx, nil)
}

func (sCtl *GoodsController) UpdateCartGoods(ctx echo.Context) error {
	userID, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "授权失败", err)
	}
	comID := echoapp_util.GetCtxComId(ctx)
	var goosdItem *echoapp.CartGoodsItem
	if err := ctx.Bind(&goosdItem); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}
	if goosdItem.Num > 10 {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "相同商品购物车支持最多添加10个商品,请分批次购买", err)
	}
	if goosdItem.Num < 0 {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}
	if err := sCtl.goodsSvr.UpdateCartGoods(comID, uint(userID), goosdItem); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "保存失败", err)
	}
	return sCtl.Success(ctx, nil)
}

func (sCtl *GoodsController) ClearCart(ctx echo.Context) error {
	userID, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "授权失败", err)
	}
	comID := echoapp_util.GetCtxComId(ctx)
	var goosdItem *echoapp.CartGoodsItem
	if err := ctx.Bind(&goosdItem); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}
	if goosdItem.Num > 10 {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "相同商品购物车支持最多添加10个商品,请分批次购买", err)
	}
	if goosdItem.Num < 0 {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}
	if err := sCtl.goodsSvr.ClearCart(comID, uint(userID)); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "保存失败", err)
	}
	return sCtl.Success(ctx, nil)
}
