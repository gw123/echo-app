package controllers

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"strconv"
)

type CommentController struct {
	echoapp.BaseController
	commentSvc echoapp.CommentService
}

func NewCommentController(commentSvc echoapp.CommentService) *CommentController {
	return &CommentController{
		commentSvc: commentSvc,
	}
}

type CommentOption struct {
	Content   string `json:"content"`
	PId       int    `json:"pid"`
	CommentId int    `json:"comment_id"`
	UserId    int    `json:"user_id"`
}

func (cmtCtrl *CommentController) SaveComment(ctx echo.Context) error {
	comment := &echoapp.Comment{}
	if err := ctx.Bind(comment); err != nil {
		return cmtCtrl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}

	if comment.Pid == 0 {
		if comment.OrderNo == "" {
			return cmtCtrl.Fail(ctx, echoapp.CodeArgument, "参数错误", errors.New("缺少order_no"))
		}

		flag, err := cmtCtrl.commentSvc.IsOrderNoExist(comment.OrderNo)
		if err != nil {
			return cmtCtrl.Fail(ctx, echoapp.CodeDBError, echoapp.ErrDb.Error(), err)
		}
		if !flag {
			return cmtCtrl.Fail(ctx, echoapp.CodeArgument, "订单不存在", errors.New("订单不存在"))
		}
	}

	userId, _ := echoapp_util.GetCtxtUserId(ctx)
	comment.UserId = userId
	if err := cmtCtrl.commentSvc.CreateComment(comment); err != nil {
		return cmtCtrl.Fail(ctx, echoapp.CodeNotFound, echoapp.ErrDb.Error(), err)
	}
	return cmtCtrl.Success(ctx, comment)
}

func (cmtCtrl *CommentController) GetCommentList(ctx echo.Context) error {
	goodsIdint, err := strconv.ParseInt(ctx.QueryParam("goods_id"), 10, 64)
	if err != nil || goodsIdint == 0 {
		return cmtCtrl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	commentlist, err := cmtCtrl.commentSvc.GetCommentList(goodsIdint, lastId, limit)
	if err != nil {
		return cmtCtrl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrDb.Error(), err)
	}
	return cmtCtrl.Success(ctx, commentlist)
}

func (cmtCtrl *CommentController) GetSubCommentList(ctx echo.Context) error {
	commentId, err := strconv.ParseInt(ctx.QueryParam("id"), 10, 64)
	if err != nil || commentId == 0 {
		return cmtCtrl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}
	comment, err := cmtCtrl.commentSvc.GetCommentById(commentId)
	if err != nil || commentId == 0 {
		return cmtCtrl.Fail(ctx, echoapp.CodeArgument,
			echoapp.ErrDb.Error(),
			errors.Wrapf(err, "GetCommentById id:%d", commentId))
	}

	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	commentlist, err := cmtCtrl.commentSvc.GetSubCommentList(commentId, lastId, limit)
	if err != nil {
		return cmtCtrl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrDb.Error(), err)
	}

	return cmtCtrl.Success(ctx, map[string]interface{}{
		"comment":        comment,
		"subCommentList": commentlist,
	})
}

func (cmtCtrl *CommentController) ThumbUpComment(ctx echo.Context) error {
	commentId := ctx.QueryParam("id")
	commentIdint, _ := strconv.ParseInt(commentId, 10, 64)
	if err := cmtCtrl.commentSvc.ThumbUpComment(commentIdint); err != nil {
		return cmtCtrl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	return cmtCtrl.Success(ctx, nil)
}

func (cmtCtrl *CommentController) GetGoodsCommentNum(ctx echo.Context) error {
	goodsIdint, err := strconv.ParseInt(ctx.QueryParam("goods_id"), 10, 64)
	if err != nil || goodsIdint == 0 {
		return cmtCtrl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}
	num, err := cmtCtrl.commentSvc.GetGoodsCommentNum(goodsIdint)
	if err != nil {
		return cmtCtrl.Fail(ctx, echoapp.CodeDBError, "查询失败", err)
	}
	return cmtCtrl.Success(ctx, map[string]interface{}{
		"num": num,
	})
}
