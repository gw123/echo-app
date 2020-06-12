package echoapp

import (
	"time"
)

type Comment struct {
	ID        int64     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	ComId     int        `json:"com_id"`
	ShopId    int        `json:"shop_id" gorm:"not null"`
	UserId    int64      `json:"user_id" form:"user_id" gorm:"not null"`
	GoodsId   int64      `json:"goods_id" gorm:"not null"`
	PId       int        `json:"pid"`
	OrderNo   string     `json:"order_no" grom:"not null"`
	SellerId  int        `json:"seller_id" gorm:"not null"`
	StaffId   int        `json:"staff_id" gorm:"not null"`
	Health    int        `json:"health" gorm:"not null"`
	Good      int        `json:"good" gorm:"not null"`
	Staff     int        `json:"staff" gorm:"not null"`
	Express   int        `json:"express" `
	Ups       int        `json:"ups" grom:"not null"`
	Covers    string     `json:"covers" gorm:"size:1024"`
	//Covers    []string   `json:"covers" gorm:"size:1024"`
	Content string `json:"content" form:"content" gorm:"size:256"`
	Source  string `json:"source"`
}

type ImageOption struct {
	ID     int
	status string
}

type CommentService interface {
	CreateComment(comment *Comment) error
	SaveComment(comment *Comment) error

	GetCommentList(goodsId int64, comId, limit int) ([]*Comment, error)
	DeleteComment(comment *Comment) error
	ThumbUpComment(commentId int64) error
	RankCommentByUp(amount int, time time.Time) error
}
