package services

import (
	"time"

	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type CommentService struct {
	db *gorm.DB
	//redis *redis.Client
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{
		db: db,
		//redis: redis,
	}
}

func (cmtSvc *CommentService) SaveComment(comment *echoapp.Comment) error {
	return cmtSvc.db.Save(comment).Error
}

func (cmtSvc *CommentService) CreateComment(comment *echoapp.Comment) error {
	//cmtSvc.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.Comment{})
	if comment.PId > 0 {
		pComment, err := cmtSvc.GetCommentById(comment.PId)
		if err != nil {
			return errors.Wrapf(err, "comment id %d not exist", comment.PId)
		}
		pComment.ReplyNum += 1
		if err = cmtSvc.SaveComment(pComment); err != nil {
			return errors.Wrapf(err, "comment save err id:%d", comment.PId)
		}
		comment.GoodsId = pComment.GoodsId
		comment.ComId = pComment.ComId
	}
	return cmtSvc.db.Create(comment).Error
}

func (cmtSvc *CommentService) GetCommentList(goodsId int64, lastId uint, limit int) ([]*echoapp.Comment, error) {
	commentList := []*echoapp.Comment{}
	query := cmtSvc.db.Where("goods_id=?", goodsId)
	if lastId > 0 {
		query = query.Where("id < ?", lastId)
	}
	res := query.Order("id desc").Where("pid = 0").Limit(limit).Find(&commentList)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "commentservice->GetCommentCommentList")
	}
	for _, item := range commentList {
		if item.Avatar == "" {
			item.Avatar = "http://img.xytschool.com/a/avatar.jpg"
		}
		if item.ReplyNum > 0 {
			cmtSvc.db.Select("id,content,created_at,avatar,nickname").
				Where("pid = ?", item.ID).Limit(2).Find(&item.ReplyList)
		}
	}
	return commentList, nil
}

func (cmtSvc *CommentService) GetGoodsCommentNum(goodsId int64) (int, error) {
	var total int
	if err := cmtSvc.db.Table("comments").Where("goods_id=?", goodsId).Count(&total).Error; err != nil {
		return 0, errors.Wrap(err, "getGoodsCommentNum")
	}
	return total, nil
}

func (cmtSvc *CommentService) GetSubCommentList(commentId int64, lastId uint, limit int) ([]*echoapp.Comment, error) {
	var commentList []*echoapp.Comment
	query := cmtSvc.db.Where("pid =?", commentId)
	if lastId > 0 {
		query = query.Where("id < ?", lastId)
	}
	res := query.Order("id desc").Limit(limit).Find(&commentList)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "commentservice->GetSubCommentCommentList")
	}
	for _, item := range commentList {
		if item.Avatar == "" {
			item.Avatar = "http://img.xytschool.com/a/avatar.jpg"
		}
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
	comment.UpNum = comment.UpNum + 1
	return cmtSvc.SaveComment(comment)
}
func (cmtSvc *CommentService) RankCommentByUp(amount int, time time.Time) error {
	panic(" a")
}

func (cmtSvc *CommentService) GetCommentById(id int64) (*echoapp.Comment, error) {
	comment := echoapp.Comment{}
	if err := cmtSvc.db.Where("id = ?", id).First(&comment).Error; err != nil {
		return nil, err
	}
	if comment.Avatar == "" {
		comment.Avatar = "http://img.xytschool.com/a/avatar.jpg"
	}
	return &comment, nil
}

func (cmtSvc *CommentService) IsOrderNoExist(orderNo string) (bool, error) {
	var count int
	if err := cmtSvc.db.Table("orders").
		Where("order_no = ?", orderNo).
		Count(&count).Error; err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
