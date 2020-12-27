package jobs

import (
	"strings"

	"github.com/gw123/echo-app/jobs"
	"github.com/pkg/errors"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/gw123/glog"
	"github.com/gw123/gworker"
	"github.com/spf13/cobra"
)

type SendSMSJobber struct {
	jobs.SendSms
}

func (s *SendSMSJobber) Handle() error {
	//发送短信
	glog.Info("发送门票短信通知 phone:" + strings.Join(s.PhoneNumbers, ","))
	sms := app.MustGetSmsService()
	if err := sms.SendMessage(&s.SendMessageOptions); err != nil {
		glog.Error(err.Error())
		return errors.Wrap(err, "smsService.SendMessage")
	}
	return nil
}

var SendSMSCmd = &cobra.Command{
	Use:   "send-sms",
	Short: "发送短信任务",
	Long:  `发送短信任务`,
	Run: func(cmd *cobra.Command, args []string) {
		model := &SendSMSJobber{}
		opt := echoapp.ConfigOpts.Job
		taskManager, err := gworker.NewConsumer(opt, model)
		if err != nil {
			glog.Errorf("NewTaskManager : %s", err.Error())
			return
		}
		taskManager.StartWork("xyt", 1)
	},
}
