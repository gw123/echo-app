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
	"os"
	"os/signal"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/services"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

func doCheckTicketWorker() {
	conn, err := amqp.Dial(echoapp.ConfigOpts.MQMap["ticket"].Url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	defer ch.Close()
	echoapp_util.DefaultLogger().Info(echoapp.ConfigOpts.TongchengConfig)
	tongchengSvr := services.NewTongchengService(echoapp.ConfigOpts.TongchengConfig)
	msgs, err := ch.Consume(
		"ticket-check", // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	for msg := range msgs {
		echoapp_util.DefaultLogger().Infof("Received a ticket-check message: %s", string(msg.Body))
		job := &echoapp.CheckTicketJob{}
		if err := json.Unmarshal(msg.Body, job); err != nil {
			echoapp_util.DefaultLogger().Errorf("Message Unmarshal: %s", err.Error())
			msg.Ack(false)
			continue
		}
		if err := tongchengSvr.CheckTicket(job); err != nil {
			echoapp_util.DefaultLogger().Errorf("CheckTicket: %s", err.Error())
			msg.Ack(false)
			continue
		}
		msg.Ack(false)
	}
}

func doSyncPartnerCodeWorker() {
	conn, err := amqp.Dial(echoapp.ConfigOpts.MQMap["ticket"].Url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	defer ch.Close()
	echoapp_util.DefaultLogger().Info(echoapp.ConfigOpts.TongchengConfig)
	tongchengSvr := services.NewTongchengService(echoapp.ConfigOpts.TongchengConfig)
	msgs, err := ch.Consume(
		"ticket-sync-code", // queue
		"",                 // consumer
		false,              // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)

	for msg := range msgs {
		echoapp_util.DefaultLogger().Infof("Received a ticket-sync-code message: %s", string(msg.Body))
		job := &echoapp.SyncPartnerCodeJob{}
		if err := json.Unmarshal(msg.Body, job); err != nil {
			echoapp_util.DefaultLogger().Errorf("Message Unmarshal: %s", err.Error())
			msg.Ack(false)
			continue
		}
		if err := tongchengSvr.SyncPartnerCode(job); err != nil {
			echoapp_util.DefaultLogger().Errorf("CheckTicket: %s", err.Error())
			msg.Ack(false)
			continue
		}
		msg.Ack(false)
	}
}

func startCheckTicketDaemon() {
	echoapp_util.DefaultLogger().Infof("开始消息推送服务验证")
	go doCheckTicketWorker()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}

func startSyncPartnerCodeDaemon() {
	echoapp_util.DefaultLogger().Infof("开始同步门票校验码到同程")
	go doSyncPartnerCodeWorker()
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

var syncPartnerCodeCmd = &cobra.Command{
	Use:   "sync-partner-code",
	Short: "同步门票核验码",
	Long:  `同步门票核验吗到同程`,
	Run: func(cmd *cobra.Command, args []string) {
		startSyncPartnerCodeDaemon()
	},
}

func init() {
	RootCmd.AddCommand(checkTicketDaemonCmd)
	RootCmd.AddCommand(syncPartnerCodeCmd)
}
