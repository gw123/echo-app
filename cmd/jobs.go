package cmd

import (
	"github.com/gw123/echo-app/cmd/crontabs"
	_ "github.com/gw123/echo-app/cmd/crontabs"
	"github.com/gw123/echo-app/cmd/jobs"
)

func init() {
	RootCmd.AddCommand(jobs.OrderCreateCmd)
	RootCmd.AddCommand(jobs.UserScoreChangeCmd)
	RootCmd.AddCommand(jobs.OrderPaidCmd)
	RootCmd.AddCommand(jobs.SyncPartnerCodeCmd)
	RootCmd.AddCommand(jobs.CheckTicketDaemonCmd)
	//
	RootCmd.AddCommand(crontabs.UpdateCacheCmd)
	RootCmd.AddCommand(crontabs.ReportTicketCmd)
}
