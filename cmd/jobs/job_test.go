package jobs

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/gw123/echo-app/jobs"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
)

func TestSendOrderPaidTplMsg(t *testing.T) {
	//echoapp.InitConfig("/Users/mac/code/go/src/github.com/gw123/echo-app/config.yaml")
	msg := &echoapp.TplMsgCreateTicket{
		BaseTemplateMessage: echoapp.BaseTemplateMessage{
			ComID:  14,
			Openid: "oNrV6w92st6gWbXxySDiohmC2KtM",
			Remark: "点击查看电子码",
			First:  "您的门票已经购买成功,点击本条消息使用门票",
		},
		UserName:   "gw123",
		OrderNO:    "12312312123123",
		TicketName: "故宫门票",
		Num:        1,
		CreatedAt:  time.Now(),
		Amount:     1.02,
	}

	wechat := app.MustGetWechatService()
	_, err := wechat.SendTplMessage(context.Background(), msg)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success")
}

func TestSendMsg(t *testing.T) {
	pusher := app.MustGetJopPusherService()
	type Params struct {
		Username string `json:"username"`
		Source   string `json:"source"`
		Code     string `json:"code"`
	}
	params := &Params{
		Username: "gw123",
		Source:   "xcx",
		Code:     "20201001",
	}
	data, _ := json.Marshal(params)
	smsJob := jobs.SendSms{
		SendMessageOptions: echoapp.SendMessageOptions{
			ComId:         14,
			PhoneNumbers:  []string{"18618184632"},
			Type:          "ticketCode",
			TemplateParam: string(data),
		},
	}
	pusher.PostJob(context.Background(), &smsJob)
}

func TestTicketeSyncCode(t *testing.T) {
	//echoapp.InitConfig("/Users/mac/code/go/src/github.com/gw123/echo-app/config.yaml")
	pusher := app.MustGetJopPusherService()
	job := jobs.TicketSyncCode{
		ComID: 14,
		Body: &echoapp.SyncPartnerCodeRequestBody{
			Tickets:        1,
			OrderSerialId:  "sz5f5c3276212379240398007",
			PartnerOrderId: "x123123",
			PartnerCode:    "2131231231",
		},
	}
	pusher.PostJob(context.Background(), &job)
}
