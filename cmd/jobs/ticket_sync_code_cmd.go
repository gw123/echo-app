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
	"github.com/gw123/glog"
	"github.com/gw123/gworker"

	"github.com/gw123/echo-app/app"

	"github.com/gw123/echo-app/jobs"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/spf13/cobra"
)

type TicketSyncCodeJobber struct {
	jobs.TicketSyncCode
}

func (t *TicketSyncCodeJobber) Handle() error {
	tc := app.GetTongchengService()
	err := tc.SyncPartnerCode(t.ComID, t.Body)

	echoapp_util.DefaultLogger().Infof("Received a TicketSyncCodeJobber message: %s", t.Body.PartnerOrderId)

	if err != nil {
		echoapp_util.DefaultLogger().Errorf("TicketSyncCodeJobber: %s", err.Error())
		return err
	}
	return nil
}

var TicketSyncCodeCmd = &cobra.Command{
	Use:   "ticket-sync-code",
	Short: "同步门票核验码",
	Long:  `同步门票核验吗到同程`,
	Run: func(cmd *cobra.Command, args []string) {
		model := &TicketSyncCodeJobber{}
		opt := echoapp.ConfigOpts.Job
		taskManager, err := gworker.NewConsumer(opt, model)
		if err != nil {
			glog.Errorf("NewTaskManager : %s", err.Error())
			return
		}
		taskManager.StartWork("xyt", 1)
	},
}
