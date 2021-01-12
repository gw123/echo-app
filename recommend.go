package echoapp

import (
	"time"
)

type EvaluationOption struct {
	UserId          int64
	GoodsId         int64
	Attribute1Score int `json:"attribute1_score"`
	Attribute2Score int `json:"attribute2_score"`
	Attribute3Score int `json:"attribute3_score"`
	Attribute4Score int `json:"attribute4_score"`
	TotalScore      int64
}

type UserCommentArray []EvaluationOption

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
	AttrituteScore  float64
}

type RecommendService interface {
	GetEvluationByUserId(userId int64) (*EvaluationOption, error)
	//GetGoodsUserComprehensivescore( int64 )(error)
}
