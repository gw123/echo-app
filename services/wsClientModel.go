package services

import (
	"encoding/json"
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/gw123/glog"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"io"
	"sync"
	"time"
)

type WsService struct {
	webSocketClientMap map[string]*WsClientModel
	UserIDClientMap    map[uint][]string
	mutex              sync.Mutex
	closeChan          chan string
}

func NewWsService() *WsService {
	return &WsService{
		webSocketClientMap: map[string]*WsClientModel{},
		closeChan:          make(chan string, 1024),
		UserIDClientMap:    make(map[uint][]string),
	}
}

//当有连接关闭时候清理
func (w *WsService) handelClose() {
	for {
		clientId := <-w.closeChan
		glog.Info(clientId)
		//避免删除的时候出现添加的情况
		w.mutex.Lock()
		client := w.webSocketClientMap[clientId]
		client.Close()
		userId := client.userID
		if userId != 0 {
			userClientList := w.UserIDClientMap[userId]
			for index, id := range userClientList {
				if id == clientId {
					w.UserIDClientMap[userId] = append(userClientList[0:index], userClientList[index:]...)
					break
				}
			}
		}
		delete(w.webSocketClientMap, clientId)
		glog.Errorf("删除clientId : %s", clientId)
		w.mutex.Unlock()
	}
}

func (w WsService) SendWsClientEvent(event echoapp.WsEvent, clientID string) error {
	client, ok := w.webSocketClientMap[clientID]
	if !ok {
		return errors.New("client not exist")
	}
	glog.Warn("ClientID:" + client.clientID)
	if err := client.SendEvent(event); err != nil {
		return errors.Wrap(err, "WsService->SendClientMessage")
	}
	return nil
}

//通过用户名来送消息
func (w WsService) SendWsClientEventByUserID(event echoapp.WsEvent, userID uint) error {
	clientIDList, ok := w.UserIDClientMap[userID]
	if !ok {
		return errors.New("client not exist")
	}
	pass := false
	for _, clientId := range clientIDList {
		if err := w.SendWsClientEvent(event, clientId); err != nil {
			glog.Errorf("SendWsClientEventByUserID userId: %d ,eventType: %s, err:%s", userID, event.GetMsgType(), err)
			continue
		} else {
			pass = true
		}
	}
	if pass == false {
		return errors.New("消息发送失败")
	}
	return nil
}

func (w WsService) AddWsClient(ctx echo.Context, userId uint, clientID string) error {
	if !ctx.IsWebSocket() {
		return errors.New("请求类型错误")
	}
	websocket.Handler(func(ws *websocket.Conn) {
		client, ok := w.webSocketClientMap[clientID]
		if ok {
			echoapp_util.ExtractEntry(ctx).Info("新的同token到来")
			client.Close()
		}
		client = NewWsClientModel(ctx, ws, clientID, userId, w.closeChan)
		w.mutex.Lock()
		w.webSocketClientMap[clientID] = client
		if userId != 0 {
			w.UserIDClientMap[userId] = append(w.UserIDClientMap[userId], clientID)
		}
		w.mutex.Unlock()

		//event := echoapp.NewWsEventNotify("测试通知", "index", "xxx", "")
		//if err := w.SendWsClientEvent(event, clientID); err != nil {
		//	glog.Error(err.Error())
		//}
		client.Run()
	}).ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}

type WsClientModel struct {
	buffer    []byte
	webSocket *websocket.Conn
	Mutex     sync.Mutex
	Token     string
	userID    uint
	clientID  string
	closeChan chan string
	runFlag   bool
	ctx       echo.Context
}

func NewWsClientModel(ctx echo.Context, conn *websocket.Conn, ClientID string, userID uint, closeChan chan string) *WsClientModel {
	this := new(WsClientModel)
	this.webSocket = conn
	this.runFlag = true
	this.buffer = make([]byte, 1024*64)
	this.closeChan = closeChan
	this.ctx = ctx
	this.clientID = ClientID
	this.userID = userID
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
	go func() {
		for ws.runFlag {
			time.Sleep(time.Second * 10)
			if err := ws.SendEvent(echoapp.NewWsEventPing()); err == io.EOF {
				ws.runFlag = false
				ws.closeChan <- ws.clientID
				glog.Warnf("ws连接断开%s userId:%d", ws.clientID, ws.userID)
				break
			}
		}
	}()

	for ws.runFlag {
		eventName, body, err := ws.ReadMsg()
		//echoapp_util.ExtractEntry(ws.ctx).Infof("eventName:%s, %s, %s", eventName, string(body), err)
		if err == io.EOF {
			//通知ws manager 清理连接
			ws.runFlag = false
			ws.closeChan <- ws.clientID
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

func (ws *WsClientModel) DealMsg(eventName string, body []byte) error {
	switch eventName {
	case echoapp.WsEventTypeCmdResult:
		result := &echoapp.WsEventCmdResult{}
		err := json.Unmarshal(body, result)
		if err != nil {
			return errors.Wrap(err, "DealMsg json.Unmarshal")
		}
		echoapp_util.ExtractEntry(ws.ctx).Infof("EventType: %s; Status: %d, result: %s", eventName, result.Status, result.Result)
	case echoapp.WsEventTypeLog:
	//echoapp_util.ExtractEntry(ws.ctx).Infof("EventType:%s;  Payload:%s", eventName, string(body))
	case echoapp.WsEventTypePing:
	case echoapp.WsEventTypePong:
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

//通知任务
//const NotifyTaskName = "notify"
//
//type NotifyTask struct {
//	TaskName string `json:"task_name"`
//	UserID   uint   `json:"user_id"`
//	ClientID string `json:"client_id"`
//	echoapp.WsEventNotify
//}
//
//func NewNotifyTask(title, position, content, link string) *NotifyTask {
//	return &NotifyTask{
//		TaskName:      NotifyTaskName,
//		WsEventNotify: *echoapp.NewWsEventNotify(title, position, content, link),
//	}
//}
//
//func (S *NotifyTask) GetTaskName() string {
//	return NotifyTaskName
//}
//
//func (S *NotifyTask) ToJson() string {
//	data, err := json.Marshal(S)
//	if err != nil {
//		return ""
//	}
//	return string(data)
//}
//
//func (*NotifyTask) GetHandleFun(ws WsService) interface{} {
//	return func(data string) error {
//		newTask := new(NotifyTask)
//		err := json.Unmarshal([]byte(data), newTask)
//		if err != nil {
//			return errors.Wrap(err, "json.Unmarshal")
//		}
//
//		if newTask.UserID != 0 {
//			app.App.WsSvr.SendWsClientEventByUserID(newTask, newTask.UserID)
//		} else {
//			app.App.WsSvr.SendWsClientEvent(newTask, newTask.ClientID)
//		}
//
//		return nil
//	}
//}