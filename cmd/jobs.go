package cmd

import (
	"os"
	"os/signal"

	"github.com/gw123/echo-app/cmd/crontabs"
	_ "github.com/gw123/echo-app/cmd/crontabs"
	"github.com/gw123/echo-app/cmd/jobs"
)

func WaitJobRunOver(job func()) {
	quit := make(chan os.Signal, 1)
	go job()
	signal.Notify(quit, os.Interrupt)
	quit <- os.Interrupt
	<-quit
}

func init() {
	// 消息队列任务
	RootCmd.AddCommand(jobs.OrderCreateCmd)
	RootCmd.AddCommand(jobs.UserScoreChangeCmd)
	RootCmd.AddCommand(jobs.OrderPaidCmd)
	RootCmd.AddCommand(jobs.TicketSyncCodeCmd)
	RootCmd.AddCommand(jobs.CheckTicketDaemonCmd)
	RootCmd.AddCommand(jobs.SendSMSCmd)

	// 定时任务
	RootCmd.AddCommand(crontabs.UpdateCacheCmd)
	RootCmd.AddCommand(crontabs.ReportTicketCmd)
	RootCmd.AddCommand(crontabs.ReportBookingPassengerFlowCmd)
}
