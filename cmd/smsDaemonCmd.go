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
<<<<<<< HEAD
	"encoding/json"
=======
	"os"
	"os/signal"

>>>>>>> feature/2020-04-26-watchercmd
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/spf13/cobra"
<<<<<<< HEAD
	"github.com/streadway/amqp"
	"os"
	"os/signal"
=======
>>>>>>> feature/2020-04-26-watchercmd
)

func doSMSWorker() {
	conn, err := amqp.Dial(echoapp.ConfigOpts.MQMap["sms"].Url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	defer ch.Close()

	smsSvr := app.MustGetSmsService()
	msgs, err := ch.Consume(
		"send-sms", // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	go func() {
		for msg := range msgs {
			echoapp_util.DefaultLogger().Infof("Received a message: %s", string(msg.Body))
			job := &echoapp.SendMessageJob{}
			if err := json.Unmarshal(msg.Body, job); err != nil {
				echoapp_util.DefaultLogger().Errorf("Message Unmarshal: %s", err.Error())
				continue
			}

			if err := smsSvr.SendMessage(&job.SendMessageOptions); err != nil {
				echoapp_util.DefaultLogger().Errorf("sendSMS: %s", err.Error())
				continue
			}
			msg.Ack(true)
		}
	}()
}
func startSmsDaemon() {
	echoapp_util.DefaultLogger().Infof("开启短信服务")
	doMqWorker()
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
