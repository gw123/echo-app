package echoapp

import (
	"encoding/json"
	"time"
)

type Comment struct {
	ID                     int64     `gorm:"primary_key" json:"id"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time
	DeletedAt              *time.Time `sql:"index"`
	ComId                  int        `json:"com_id"`
	ShopId                 int        `json:"shop_id" gorm:"not null"`
	UserId                 int64      `json:"user_id" form:"user_id" gorm:"not null"`
	GoodsId                int64      `json:"goods_id" gorm:"not null"`
	PId                    int64      `json:"pid"`
	OrderNo                string     `json:"order_no" grom:"not null"`
	SellerId               int        `json:"seller_id" gorm:"not null"`
	StaffId                int        `json:"staff_id" gorm:"not null"`
	Health                 int        `json:"health" gorm:"not null"`
	Good                   int        `json:"good" gorm:"not null"`
	Staff                  int        `json:"staff" gorm:"not null"`
	Express                int        `json:"express" `
	UpNum                  int        `json:"up_num" gorm:"not null"`
	ReplyNum               int        `json:"reply_num" gorm:"not null"`
	CoversStr              string     `json:"-" gorm:"column:covers;size:1024"`
	Covers                 []string   `json:"covers" gorm:"-"`
	Content                string     `json:"content" form:"content" gorm:"size:256"`
	Source                 string     `json:"source"`
	Avatar                 string     `json:"avatar"`
	NickName               string     `json:"nickname"`
	ReplyList              []*Comment `json:"reply_list" gorm:"-"`
	UserComprehensiveScore float64    `json:"score" gorm:"column:score"`
}

func (c *Comment) BeforeCreate() (err error) {
	str, err := json.Marshal(c.Covers)
	if err != nil {
		return err
	}
	c.CoversStr = string(str)
	return
}
func (c *Comment) AfterFind() error {
	if err := json.Unmarshal([]byte(c.CoversStr), &c.Covers); err != nil {
		return err
	}
	return nil
}

type CommentService interface {
	CreateComment(comment *Comment) error
	SaveComment(comment *Comment) error

	GetCommentList(goodsId int64, comId, limit int) ([]*Comment, error)
	DeleteComment(comment *Comment) error
	ThumbUpComment(commentId int64) error
	RankCommentByUp(amount int, time time.Time) error

	GetCommentById(id int64) (*Comment, error)
	IsOrderNoExist(orderNo string) (bool, error)
	GetGoodsCommentNum(goodsId int64) (int, error)
	GetSubCommentList(id int64, lastId int, limit int) ([]*Comment, error)
}
