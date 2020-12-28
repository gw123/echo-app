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
	"github.com/gw123/echo-app/app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/gw123/glog"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"time"
)

var updateMethod = ""

func updateUserCache() {
	echoapp_util.DefaultLogger().Info("开启更新user缓存服务")
	//echoapp_util.DefaultLogger().Infof("%+v", echoapp.ConfigOpts)
	userDb := app.MustGetDb("user")
	//cache := app.MustGetRedis("")
	usrSvr := app.MustGetUserService()
	userList := []*echoapp.User{}
	var currentMaxId int64 = 0
	var size int64 = 100

	for {
		echoapp_util.DefaultJsonLogger().Errorf("更新缓存 from %d to %d", currentMaxId, currentMaxId+size)
		//roles := []echoapp.Role{}
		if err := userDb.
			Where("id > ?", currentMaxId).
			Where("com_id != 0").
			Order("id asc").
			Limit(size).
			Find(&userList).Error; err != nil {
			break
		}

		if len(userList) == 0 {
			echoapp_util.DefaultJsonLogger().Info("更新结束")
			break
		}
		//更新用户角色
		for _, user := range userList {
			if err := usrSvr.UpdateCachedUser(user); err != nil {
				echoapp_util.DefaultJsonLogger().WithError(err).Error("更新用户缓存失败")
			}
			currentMaxId = user.Id
		}
		time.Sleep(time.Second * 2)
	}

}

func updateCompanyCache() {
	echoapp_util.DefaultLogger().Info("开启更新com缓存服务")
	companySvr := app.MustGetCompanyService()
	var currentMaxId uint = 0
	var limit uint = 100

	for {
		echoapp_util.DefaultJsonLogger().Errorf("更新缓存 from %d to %d", currentMaxId, currentMaxId+limit)
		list, err := companySvr.GetCompanyList(currentMaxId, int(limit))
		if err != nil {
			echoapp_util.DefaultJsonLogger().WithError(err).Errorf("GetCompnayList:%s", err.Error())
		}

		if len(list) == 0 {
			echoapp_util.DefaultJsonLogger().Info("更新结束")
			break
		}

		for _, company := range list {
			companyDetail, err := companySvr.GetCompanyById(company.Id)
			if err != nil {
				glog.Error(err.Error())
				continue
			}
			if err := companySvr.UpdateCachedCompany(companyDetail); err != nil {
				echoapp_util.DefaultJsonLogger().WithError(err).Error("更新用户缓存失败")
			}
			currentMaxId = company.Id
		}
		time.Sleep(time.Second * 2)
	}

}

func updateCouponCache() {
	echoapp_util.DefaultLogger().Info("开启更新coupon缓存")
	companySvr := app.MustGetCompanyService()
	actSvr := app.MustGetActivityService()
	var currentMaxId uint = 0
	var limit uint = 100

	for {
		echoapp_util.DefaultJsonLogger().Errorf("更新缓存 from %d to %d", currentMaxId, currentMaxId+limit)
		list, err := companySvr.GetCompanyList(currentMaxId, int(limit))
		if err != nil {
			echoapp_util.DefaultJsonLogger().WithError(err).Errorf("GetCompnayList:%s", err.Error())
		}

		if len(list) == 0 {
			echoapp_util.DefaultJsonLogger().Info("更新结束")
			break
		}

		for _, company := range list {
			_, err := actSvr.UpdateCachedCouponsByComId(company.Id, 0)
			if err != nil {
				echoapp_util.DefaultJsonLogger().WithError(err).Errorf("UpdateCachedCouponsByComId:%s", err.Error())
			}
			currentMaxId = company.Id
		}
		time.Sleep(time.Second * 2)
	}

}

func updateCacheByMq() {
	conn, err := amqp.Dial(echoapp.ConfigOpts.MQMap["sms"].Url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	defer ch.Close()

	usrSvr := app.MustGetUserService()
	comSvr := app.MustGetCompanyService()

	msgs, err := ch.Consume(
		"update-cache", // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	for msg := range msgs {
		echoapp_util.DefaultLogger().Infof("Received a message: %s", string(msg.Body))
		job := &echoapp.UpdateCacheData{}
		if err := json.Unmarshal(msg.Body, job); err != nil {
			echoapp_util.DefaultLogger().Errorf("Message Unmarshal: %s", err.Error())
			msg.Ack(true)
			continue
		}
		func() {
			switch job.Type {
			case echoapp.RedisUserKey:
				user, err := usrSvr.GetUserById(job.UserId)
				if err != nil {
					echoapp_util.DefaultLogger().Errorf("getUserById: %d,err:%s", job.UserId, err.Error())
					break
				}
				if err := usrSvr.UpdateCachedUser(user); err != nil {
					echoapp_util.DefaultLogger().Errorf("updateCacheUser: %d,err:%s", job.UserId, err.Error())
				}

			case echoapp.RedisCompanyKey:
				com, err := comSvr.GetCompanyById(job.ComId)
				if err != nil {
					echoapp_util.DefaultLogger().Errorf("getCompanyById: %d,err:%s", job.ComId, err.Error())
					break
				}
				if err := comSvr.UpdateCachedCompany(com); err != nil {
					echoapp_util.DefaultLogger().Errorf("updateCacheCompany: %d,err:%s", job.UserId, err.Error())
				}
			}
		}()
		msg.Ack(true)
	}
}

// serverCmd represents the server command
var updateCacheCmd = &cobra.Command{
	Use:   "updateCache",
	Short: "cache 更新",
	Long:  `cache 服务`,
	Run: func(cmd *cobra.Command, args []string) {

		quit := make(chan os.Signal, 1)
		switch updateMethod {
		case "once":
			go func() {
				updateCompanyCache()
				updateCouponCache()
				updateUserCache()
				quit <- os.Interrupt
			}()
		case "mq":
			updateCacheByMq()
		}
		signal.Notify(quit, os.Interrupt)
		<-quit
	},
}

func init() {
	updateCacheCmd.PersistentFlags().StringVarP(&updateMethod, "method", "m", "once", "访问方法")
	rootCmd.AddCommand(updateCacheCmd)
}
