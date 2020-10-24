package echoapp

import "time"

const (
	ViodeStatusOnline  = 1 // 视频在线
	ViodeStatusOffline = 2 // 视频下线
)

type Video struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	ComID      uint      `json:"com_id"`
	Type       string    `json:"type"`   //视频类型 RTP/RTCP/RTSP/RTMP/MMS/HLS
	Status     int       `json:"status"` //状态  1,可以播放  2,下线
	Src        string    `json:"src"`
	Name       string    `json:"name"`
	SmallCover string    `json:"small_cover"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type VideoService interface {
	GetVideoList(comId uint, lastId uint, limit int) ([]*Video, error)
	GetVideoDetail(id uint) (*Video, error)
}
