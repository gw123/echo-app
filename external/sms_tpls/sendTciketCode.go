package sms_tpls

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/gw123/echo-app/jobs"
)

/**
尊敬的游客，您在网上订票成功,
请点击https://m.xx.com/ticket/5134233753351226
查看您的电子门票,客服电话18603868343。
*/

const SMS_Tpl_Ticket_Code = "ticketCode"

func SendTicketCode(comID uint, mobile, username, source, code string) error {
	pusher, err := app.GetJobPusher()
	if err != nil {
		return errors.Wrap(err, "SendTicketCode GetJobPusher")
	}

	type Params struct {
		Username string `json:"username"`
		Source   string `json:"source"`
		Code     string `json:"code"`
	}
	params := &Params{
		Username: username,
		Source:   source,
		Code:     code,
	}
	data, _ := json.Marshal(params)
	smsJob := jobs.SendSms{
		SendMessageOptions: echoapp.SendMessageOptions{
			ComId:         comID,
			PhoneNumbers:  []string{mobile},
			Type:          SMS_Tpl_Ticket_Code,
			TemplateParam: string(data),
		},
	}
	if err := pusher.PostJob(context.Background(), &smsJob); err != nil {
		return errors.Wrap(err, SMS_Tpl_Ticket_Code)
	}
	return nil
}
