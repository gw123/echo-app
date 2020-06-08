package services

import (
	"time"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type CommentService struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewCommentService(db *gorm.DB, redis *redis.Client) *CommentService {
	return &CommentService{
		db:    db,
		redis: redis,
	}
}
func (cmtSvc *CommentService) SaveComment(comment *echoapp.Comment) error {
	return cmtSvc.db.Save(comment).Error
}
func (cmtSvc *CommentService) CreateComment(comment *echoapp.Comment) error {
	return cmtSvc.db.Create(comment).Error
}
func (cmtSvc *CommentService) GetCommentList(goodsId int64, comId, limit int) ([]*echoapp.Comment, error) {
	commentList := []*echoapp.Comment{}
	res := cmtSvc.db.Where("com_id=? and goods_id=?", comId, goodsId).Order("created_at desc").Limit(limit).Find(&commentList)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "commentservice->GetCommentCommentList")
	}
	return commentList, nil
}
func (cmtSvc *CommentService) DeleteComment(comment *echoapp.Comment) error {
	return cmtSvc.db.Delete(comment).Error
}
func (cmtSvc *CommentService) ThumbUpComment(commentId int64) error {
	comment := &echoapp.Comment{}
	res := cmtSvc.db.Where("id=?", commentId).First(comment)
	if res.Error != nil {
		return errors.Wrap(res.Error, "commentservice->ThumbUpComment")
	}
	comment.Ups++
	return cmtSvc.SaveComment(comment)
}
func (cmtSvc *CommentService) RankCommentByUp(amount int, time time.Time) error {
	panic(" a")
}
