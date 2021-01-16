package sms_tpls

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/jobs"
)

/**
尊敬的XX,您预约的2020-01-02日09:00-12:00点的门票己经成功，
请您携带此码和门票码一起入园使用。
https://m.xx.com/appointment/5134233753351226查看您的预约码。
客服09008782.大约这样子的
*/

const SMS_Tpl_Appointment_Code = "appointment_code"

func (s *SmsTplAPiIpl) SendAppointmentCode(comID uint, mobile, username, source, date, timeRange, code string) error {
	//pusher, err := app.GetJobPusher()
	//if err != nil {
	//	return errors.Wrap(err, "SendAppointmentCode GetJobPusher")
	//}

	params := struct {
		Username  string `json:"username"`
		Source    string `json:"source"`
		Code      string `json:"code"`
		Date      string `json:"date"`
		TimeRange string `json:"timeRange"`
	}{
		Username:  username,
		Source:    source,
		Code:      code,
		Date:      date,
		TimeRange: timeRange,
	}

	data, _ := json.Marshal(params)
	smsJob := jobs.SendSms{
		SendMessageOptions: echoapp.SendMessageOptions{
			ComId:         comID,
			PhoneNumbers:  []string{mobile},
			Type:          SMS_Tpl_Appointment_Code,
			TemplateParam: string(data),
		},
	}
	if err := s.pusher.PostJob(context.Background(), &smsJob); err != nil {
		return errors.Wrap(err, SMS_Tpl_Appointment_Code)
	}
	return nil
}
