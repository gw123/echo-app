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

package cmd

import (
	"os"
	"os/signal"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/spf13/cobra"
)

func startSmsDaemon() {
	echoapp_util.DefaultLogger().Infof("%+v", echoapp.ConfigOpts)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

// serverCmd represents the server command
var smsDaemonCmd = &cobra.Command{
	Use:   "sms",
	Short: "短信",
	Long:  `短信消费者`,
	Run: func(cmd *cobra.Command, args []string) {
		startSmsDaemon()
	},
}

func init() {
	rootCmd.AddCommand(smsDaemonCmd)
}
