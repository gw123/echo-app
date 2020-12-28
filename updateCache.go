package echoapp

type UpdateCacheData struct {
	Type       string `json:"type"`
	ComId      uint    `json:"com_id"`
	UserId     int64  `json:"user_id"`
	GoodsId    int64  `json:"goods_id"`
	ArticleId  int    `json:"article_id"`
	ActivityId int    `json:"activity_id"`
}

type UpdateCacheJob struct {
	BaseMqMsg
	UpdateCacheData
}

type UpdateCacheService interface {
	Update(opt *UpdateCacheJob) error
}
