package jobs

import (
	"github.com/gw123/echo-app/jobs"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/gw123/glog"
	"github.com/gw123/gworker"
	"github.com/spf13/cobra"
)

type UserScoreChangeJobber struct {
	jobs.UserScoreChange
}

func (s *UserScoreChangeJobber) Handle() error {
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
		model := &UserScoreChangeJobber{}
		opt := echoapp.ConfigOpts.Job
		opt.DefaultQueue = model.GetName()
		taskManager, err := gworker.NewConsumer(opt, "xyt")
		if err != nil {
			glog.Errorf("NewTaskManager : %s", err.Error())
			return
		}
		taskManager.RegisterTask(model)
		taskManager.StartWork("xyt", 1)
	},
}
