// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jobs

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/gw123/echo-app/jobs"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/gw123/glog"
	"github.com/gw123/gworker"
	"github.com/spf13/cobra"
)

type TicketCheckJobber struct {
	jobs.TicketCheck
}

func (t *TicketCheckJobber) Handle() error {
	tc := app.GetTongchengService()
	err := tc.CheckTicket(t.ComID, &t.CheckTicketRequestBody)
	echoapp_util.DefaultLogger().Infof("Received a ticket-check message: %s", t.OrderSerialId)
	if err != nil {
		echoapp_util.DefaultLogger().Errorf("CheckTicket: %s", err.Error())
		return err
	}
	return nil
}

// serverCmd represents the server command
var CheckTicketDaemonCmd = &cobra.Command{
	Use:   "check-ticket",
	Short: "验票核销消息",
	Long:  `验票消费者`,
	Run: func(c *cobra.Command, args []string) {
		model := &TicketCheckJobber{}
		opt := echoapp.ConfigOpts.Job
		taskManager, err := gworker.NewConsumer(opt, model)
		if err != nil {
			glog.DefaultLogger().Errorf("NewTaskManager : %s", err.Error())
			return
		}
		taskManager.StartWork("xyt", 1)
	},
}
