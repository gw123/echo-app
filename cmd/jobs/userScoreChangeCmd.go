package jobs

import (
	"github.com/gw123/echo-app/jobs"

	"github.com/RichardKnop/machinery/v1/config"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/gw123/glog"
	"github.com/gw123/gworker"
	"github.com/spf13/cobra"
)

type UserScoreChange struct {
	jobs.SendSms
}

func (s *UserScoreChange) Handle() error {
	glog.DefaultLogger().WithField("userScore", s).Info("积分变化")
	userService := app.MustGetUserService()
	err := userService.AddScoreByUserId(s.ComId, s.UserId, s.Score, s.Source, s.SourceDetail, s.Note)
	if err != nil {
		glog.DefaultLogger().WithError(err).Errorf("AddSourceByUserId ")
		return err
	}
	return nil
}

var UserScoreChangeCmd = &cobra.Command{
	Use:   "user_score_change",
	Short: "更新用户积分",
	Long:  `更新用户积分`,
	Run: func(cmd *cobra.Command, args []string) {
		model := &UserScoreChange{}
		opt := echoapp.ConfigOpts.Job
		cfg := &config.Config{
			Broker:        opt.Broker,
			DefaultQueue:  model.GetName(),
			ResultBackend: opt.ResultBackend,
			AMQP: &config.AMQPConfig{
				Exchange:      opt.AMQP.Exchange,
				ExchangeType:  opt.AMQP.ExchangeType,
				PrefetchCount: opt.AMQP.PrefetchCount,
				AutoDelete:    opt.AMQP.AutoDelete,
			},
		}
		taskManager, err := gworker.NewConsumer(cfg, "xyt")
		if err != nil {
			glog.Errorf("NewTaskManager : %s", err.Error())
			return
		}
		taskManager.RegisterTask(&UserScoreChange{})
		taskManager.StartWork("xyt", 1)
	},
}
