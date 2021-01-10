package echoapp

import (
	"github.com/labstack/echo"
	"time"
)

const (
	WsEventTypeLog       = "log"
	WsEventTypePing      = "ping"
	WsEventTypePong      = "pong"
	WsEventTypeCmd       = "cmd"
	WsEventTypeCmdResult = "cmd_result"

	WsEventTypeNotify = "notify"
)

type WsService interface {
	AddWsClient(ctx echo.Context, userId uint, token string) error
	SendWsClientEvent(event WsEvent, token string) error
	SendWsClientEventByUserID(event WsEvent, userID uint) error
}

type WsEvent interface {
	GetMsgType() string
	GetMsgId() string
}

type WsClient interface {
	Close() error
	Run() error
	SendEvent(event WsEvent) error
}

//基础消息类型
type WsEventBase struct {
	EventId   string `json:"event_id"`
	EventType string `json:"event_type"`
	Source    string `json:"source"`
	CreatedAt int64  `json:"created_at"`
}

func (w *WsEventBase) GetMsgType() string {
	return w.EventType
}

func (w *WsEventBase) GetMsgId() string {
	return w.GetMsgId()
}

//上报客户端日志
type WsEventLog struct {
	WsEventBase
	Level    string            `json:"level"`
	Content  string            `json:"content"`
	Position string            `json:"position"`
	Tags     map[string]string `json:"tags"`
}

//ping
type WsEventPing struct {
	WsEventBase
}

func NewWsEventPing() *WsEventPing {
	return &WsEventPing{
		WsEventBase: WsEventBase{
			EventId:   time.Now().Local().String(),
			EventType: WsEventTypePing,
			Source:    "server",
			CreatedAt: time.Now().Unix(),
		},
	}
}

//命令
type WsEventCmd struct {
	WsEventBase
	Cmd       string   `json:"cmd"`
	RequestId string   `json:"request_id"`
	Params    []string `json:"params"`
}

//执行结果
type WsEventCmdResult struct {
	WsEventBase
	Status int    `json:"status"`
	Result string `json:"result"`
}

//ping
type WsEventNotify struct {
	WsEventBase
	Title    string `json:"title"`
	Position string `json:"position"`
	Content  string `json:"content"`
	Link     string `json:"link"`
}

//创建新的通知消息
func NewWsEventNotify(title, position, content, link string) *WsEventNotify {
	return &WsEventNotify{
		WsEventBase: WsEventBase{
			EventId:   time.Now().Local().String(),
			EventType: WsEventTypeNotify,
			Source:    "server",
			CreatedAt: time.Now().Unix(),
		},
		Title:    title,
		Position: position,
		Content:  content,
		Link:     link,
	}
}


