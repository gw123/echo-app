package echoapp

import (
	"github.com/labstack/echo"
)

const (
	WsEventTypeLog       = "log"
	WsEventTypePing      = "ping"
	WsEventTypeCmd       = "cmd"
	WsEventTypeCmdResult = "cmd_result"
)

type WsService interface {
	AddWsClient(ctx echo.Context, token string) error
	DelWsClient(token string) error
	SendWsClientEvent(event WsEvent, token string) error
}

type WsEvent interface {
	GetMsgType() string
	GetMsgId() string
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

type WsClient interface {
	Close() error
	Run() error
	SendEvent(event WsEvent) error
}
