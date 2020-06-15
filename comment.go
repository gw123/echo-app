package echoapp

import (
	"encoding/json"
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
	Pid       int64      `json:"pid"`
	OrderNo   string     `json:"order_no" grom:"not null"`
	SellerId  int        `json:"seller_id" gorm:"not null"`
	StaffId   int        `json:"staff_id" gorm:"not null"`
	Health    int        `json:"health" gorm:"not null"`
	Goods     int        `json:"goods" gorm:"not null"`
	Staff     int        `json:"staff" gorm:"not null"`
	Express   int        `json:"express" `
	UpNum     int        `json:"up_num" gorm:"not null"`
	ReplyNum  int        `json:"reply_num" gorm:"not null"`
	Covers    []string   `json:"covers" gorm:"-"`
	CoversStr string     `gorm:"column:covers;size:1024" json:"-"`
	Content   string     `json:"content" form:"content" gorm:"size:256"`
	Source    string     `json:"source"`
	Avatar    string     `json:"avatar"`
	Nickname  string     `json:"nickname"`
	ReplyList []*Comment `json:"reply_list" gorm:"-"`
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

type EvaluationOption struct {
	UserId          int64
	GoodsId         int64
	Attribute1Score string `json:"attribute1_score"`
	Attribute2Score string `json:"attribute2_score"`
	Attribute3Score string `json:"attribute3_score"`
	Attribute4Score string `json:"attribute4_score"`
	TotalScore      int64
}

type UserCommentArray []EvaluationOption
type ImageOption struct {
	ID     int
	status string
}

// 用户集
type UserOotions struct {
	UserId   int64
	Sex      string
	Role     string
	City     string `json:"city"`
	Score    int    `json:"score"`
	IsStaff  bool   `json:"is_staff"`
	IsVip    string `json:"is_vip"`
	VipLevel string `json:"vip_level"`
}

// 项目集(项目以goods为例)
type GoodsOptions struct {
	GoodsId    int64
	CreatedAt  time.Time
	UserId     int64   `json:"user_id"`
	TagStr     string  `josn:"tags"`
	Tags       string  `josn:"-" gorm:"-"`
	Name       string  `json:"name"`
	Price      float32 `json:"price"`
	Body       string  `json:"body"`
	RealPrice  float32 `json:"real_price"`
	GoodType   string  `json:"good_type"`
	Status     string  `json:"status"`
	SmallCover string  ` json:"small_cover"`
	Covers     string  `gorm:"type:varchar(2048)" json:"covers"`
	Pages      int     `json:"pages"`
}

// 用户-项目评分集
type UserGoodsEvaluationOptions struct {
	GoodsId         int64  `json:"goods_id"`
	UserId          int64  `json:"user_id"`
	Attribute1Score string `json:"attribute1_score"`
	Attribute2Score string `json:"attribute2_score"`
	Attribute3Score string `json:"attribute3_score"`
	Attribute4Score string `json:"attribute4_score"`
	AttrituteScore  map[string]int
}

type CommentService interface {
	CreateComment(comment *Comment) error
	SaveComment(comment *Comment) error
	//GetCommentById(id int) (*Comment, error)
	//GetCommentByTargetId(targetId int64, limit int) (*Comment, error)
	GetCommentList(goodsId int64, lastId, limit int) ([]*Comment, error)
	//UpdateComment(comment *Comment) error
	DeleteComment(comment *Comment) error
	ThumbUpComment(commentId int64) error
	RankCommentByUp(amount int, time time.Time) error
	GetCommentById(id int64) (*Comment, error)
	IsOrderNoExist(orderNo string) (bool, error)
	GetGoodsCommentNum(goodsId int64) (int, error)
	GetSubCommentList(id int64, lastId int, limit int) ([]*Comment, error)
}
