package jobs

import (
	"context"
	"testing"
	"time"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
)

func TestSendMsg(t *testing.T) {

}

func TestSendOrderPaidTplMsg(t *testing.T) {
	echoapp.InitConfig("/Users/mac/code/go/src/github.com/gw123/echo-app/config.yaml")
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
