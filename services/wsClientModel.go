package services

import (
	"encoding/json"
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"io"
	"sync"
)

type WsService struct {
	webSocketClientMap map[string]*WsClientModel
	mutex              sync.Mutex
	closeChan          chan string
}

func NewWsService() *WsService {
	return &WsService{
		webSocketClientMap: map[string]*WsClientModel{},
		closeChan:          make(chan string, 1024),
	}
}

func (w WsService) SendWsClientEvent(event echoapp.WsEvent, token string) error {
	client, ok := w.webSocketClientMap[token]
	if !ok {
		return errors.New("client not exist")
	}
	if err := client.SendEvent(event); err != nil {
		return errors.Wrap(err, "WsService->SendClientMessage")
	}
	return nil
}

func (w WsService) AddWsClient(ctx echo.Context, token string) error {
	if !ctx.IsWebSocket() {
		return errors.New("请求类型错误")
	}

	websocket.Handler(func(ws *websocket.Conn) {
		client, ok := w.webSocketClientMap[token]
		if ok {
			echoapp_util.ExtractEntry(ctx).Info("新的同名模块连接到来")
			client.Close()
		}
		client = NewWsClientModel(ctx, ws, w.closeChan)
		w.mutex.Lock()
		w.webSocketClientMap[token] = client
		w.mutex.Unlock()
		client.Run()
	}).ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

func (w WsService) DelWsClient(token string) error {
	panic("implement me")
}

type WsClientModel struct {
	buffer    []byte
	webSocket *websocket.Conn
	Mutex     sync.Mutex
	Token     string
	User      echoapp.User
	closeChan chan string
	runFlag   bool
	ctx       echo.Context
}

func NewWsClientModel(ctx echo.Context, conn *websocket.Conn, closeChan chan string) *WsClientModel {
	this := new(WsClientModel)
	this.webSocket = conn
	this.runFlag = true
	this.buffer = make([]byte, 1024*64)
	this.closeChan = closeChan
	this.ctx = ctx
	return this
}

func (ws *WsClientModel) SendEvent(event echoapp.WsEvent) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	_, err = ws.webSocket.Write(eventData)
	if err != nil {
		return err
	}
	return nil
}

func (ws *WsClientModel) Run() error {
	for ws.runFlag {
		eventName, body, err := ws.ReadMsg()
		//echoapp_util.ExtractEntry(ws.ctx).Infof("eventName:%s, %s, %s", eventName, string(body), err)
		if err == io.EOF {
			//通知ws manager 清理连接
			ws.runFlag = false
			ws.closeChan <- ws.Token
			break
		}

		if err != nil {
			echoapp_util.ExtractEntry(ws.ctx).Infof("消息解析失败:%s", err.Error())
			continue
		}

		if err = ws.DealMsg(eventName, body); err != nil {
			echoapp_util.ExtractEntry(ws.ctx).Infof("消息处理失败:%s", err.Error())
		}
	}
	return nil
}

func (ws *WsClientModel) DealMsg(eventName string, body []byte, ) error {
	switch eventName {
	case echoapp.WsEventTypeCmdResult:
		result := &echoapp.WsEventCmdResult{}
		err := json.Unmarshal(body, result)
		if err != nil {
			return errors.Wrap(err, "DealMsg json.Unmarshal")
		}
		echoapp_util.ExtractEntry(ws.ctx).Infof("EventType: %s; Status: %s, result: %s",
			eventName, result.Status, result.Result)
	case echoapp.WsEventTypeLog: fallthrough
	case echoapp.WsEventTypePing:
		echoapp_util.ExtractEntry(ws.ctx).Infof("EventType:%s;  Payload:%s", eventName, string(body))
	default:
		echoapp_util.ExtractEntry(ws.ctx).Infof("未知消息类型%s", string(body))
	}
	return nil
}

func (ws *WsClientModel) ReadMsg() (eventName string, body []byte, err error) {
	n, err := ws.webSocket.Read(ws.buffer)
	if err != nil {
		return "", nil, err
	}
	var decodeBufer []byte
	for pos, b := range ws.buffer {
		if b == '#' {
			decodeHeader := ws.buffer[:pos]
			decodeBufer = ws.buffer[pos+1 : n]
			return string(decodeHeader), decodeBufer, nil
		}
	}
	return "", nil, errors.New("协议识别失败")
}

func (ws *WsClientModel) Close() error {
	ws.runFlag = false
	if ws.webSocket.IsClientConn() || ws.webSocket.IsServerConn() {
		return ws.webSocket.Close()
	}
	return nil
}
