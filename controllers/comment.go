package controllers

import (
	"strconv"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
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

// type CommentOption struct {
// 	Content   string `json:"content"`
// 	PId       int    `json:"pid"`
// 	CommentId int    `json:"comment_id"`
// 	UserId    int    `json:"user_id"`
// }

func (cmtCtrl *CommentController) SaveComment(ctx echo.Context) error {
	comment := &echoapp.Comment{}
	if err := ctx.Bind(comment); err != nil {
		return cmtCtrl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	userId, _ := echoapp_util.GetCtxtUserId(ctx)
	comment.UserId = userId

	health := echoapp_util.TFSToFS(echoapp_util.LinguisticToTFS(comment.Health))
	good := echoapp_util.TFSToFS(echoapp_util.LinguisticToTFS(comment.Good))
	staff := echoapp_util.TFSToFS(echoapp_util.LinguisticToTFS(comment.Staff))
	//express := echoapp_util.TFSToFS(echoapp_util.LinguisticToTFS(comment.Express))

	comment.UserComprehensiveScore, _ = echoapp_util.WFGHM(1.0,
		2.0, []float64{health, good, staff}, []float64{0.3, 0.5, 0.2})

	if err := cmtCtrl.commentSvc.CreateComment(comment); err != nil {
		return cmtCtrl.Fail(ctx, echoapp.CodeNotFound, echoapp.ErrNotFoundDb.Error(), err)
	}
	return cmtCtrl.Success(ctx, comment)
}

func (cmtCtrl *CommentController) GetCommentList(ctx echo.Context) error {
	goodsId := ctx.QueryParam("goodsId")
	goodsIdint, _ := strconv.ParseInt(goodsId, 10, 64)
	limit := ctx.QueryParam("limit")
	limitint, _ := strconv.Atoi(limit)
	comId := ctx.QueryParam("com_id")
	comIdInt, _ := strconv.Atoi(comId)
	commentlist, err := cmtCtrl.commentSvc.GetCommentList(goodsIdint, comIdInt, limitint)
	if err != nil {
		return cmtCtrl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	return cmtCtrl.Success(ctx, commentlist)
}
func (cmtCtrl *CommentController) ThumbUpComment(ctx echo.Context) error {
	commentId := ctx.QueryParam("commentId")
	commentIdint, _ := strconv.ParseInt(commentId, 10, 64)
	if err := cmtCtrl.commentSvc.ThumbUpComment(commentIdint); err != nil {
		return cmtCtrl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	return cmtCtrl.Success(ctx, nil)
}
