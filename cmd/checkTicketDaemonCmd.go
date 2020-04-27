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
	"encoding/json"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/services"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
)

func doMqWorker() {
	conn, err := amqp.Dial(echoapp.ConfigOpts.MQMap["ticket"].Url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	defer ch.Close()

	tongchengSvr := services.NewTongchengService(echoapp.ConfigOpts.TongChengMap)
	msgs, err := ch.Consume(
		"check-ticket", // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	for msg := range msgs {
		echoapp_util.DefaultLogger().Infof("Received a message: %s", string(msg.Body))
		job := &echoapp.CheckTicketJob{}
		if err := json.Unmarshal(msg.Body, job); err != nil {
			echoapp_util.DefaultLogger().Errorf("Message Unmarshal: %s", err.Error())
			continue
		}
		if err := tongchengSvr.CheckTicket(job); err != nil {
			echoapp_util.DefaultLogger().Errorf("CheckTicket: %s", err.Error())
			continue
		}
		msg.Ack(true)
	}
}

func startCheckTicketDaemon() {
	echoapp_util.DefaultLogger().Infof("开始消息推送服务")
	go doMqWorker()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

// serverCmd represents the server command
var checkTicketDaemonCmd = &cobra.Command{
	Use:   "check-ticket",
	Short: "验票核销消息",
	Long:  `验票消费者`,
	Run: func(cmd *cobra.Command, args []string) {
		startCheckTicketDaemon()
	},
}

func init() {
	rootCmd.AddCommand(checkTicketDaemonCmd)
}
