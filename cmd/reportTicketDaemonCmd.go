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
	"bytes"
	"encoding/json"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"time"
)

type RealGateDay struct {
	ScenicCode string `json:"scenicCode"`
	Day        string `json:"day"`
	Count      int    `json:"count"`
}

type RealPeopleNumber struct {
	ScenicCode string `json:"scenicCode"`
	UpTime     string `json:"upTime"`
	Total      int    `json:"total"`
}

type ReportDataResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type ReportDataResponseV2 struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

type ChannelHourReportV2 struct {
	InNum      int    `json:"inNum" gorm:"column:in_num" `
	OutNum     int    `json:"outNum" gorm:"column:out_num"`
	ChannelId  string `json:"channelId" gorm:"column:channel_id"`
	RecordTime string `json:"recordTime"`
}

type ReportHourV2 struct {
	LoginName string                 `json:"loginName"`
	Pwd       string                 `json:"pwd"`
	Data      []*ChannelHourReportV2 `json:"data"`
}

func doReportHttpRequest(url, app_key string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->http.NewRequest")
	}
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set("appKey", app_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->Do")
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->ReadAll")
	}
	if res.StatusCode == http.StatusOK {
		return data, errors.Wrapf(err, "StatusCode:%d", res.StatusCode)
	}
	return data, nil
}

func reportDailyTicket() error {
	clientMap := echoapp.ConfigOpts.ReportTicketMap
	for key, options := range clientMap {
		echoapp_util.DefaultLogger().Infof("推送%s流量,com_id:%d", key, options.ComId)
		now := time.Now()
		today := now.Format("2006-01-02")
		nextDay := now.Add(time.Hour * 24).Format("2006-01-02")
		db, err := app.GetDb("shop")
		if err != nil {
			return errors.Wrap(err, "GetDb")
		}

		var count int
		if err := db.Table("tickets").
			Where("com_id = ?", options.ComId).
			Where("used_at > ? and used_at < ?", today, nextDay).
			Count(&count).Error; err != nil {
			return errors.Wrap(err, "Query")
		}

		reportDataRequest := []RealGateDay{{
			ScenicCode: options.ScenicCode,
			Day:        today,
			Count:      count,
		}}

		data, err := json.Marshal(reportDataRequest)
		if err != nil {
			return errors.Wrap(err, "json.Marshal")
		}
		echoapp_util.DefaultLogger().Info(string(data))
		url := options.BaseUrl + "/upload-data/tourist/real-gate-day"
		responseData, err := doReportHttpRequest(url, options.AppKey, data)
		if err != nil {
			return errors.Wrap(err, "doReportHttpRequest")
		}

		reportDataResponse := &ReportDataResponse{}
		if err = json.Unmarshal(responseData, reportDataResponse); err != nil {
			return errors.Wrap(err, "doReportHttpRequest")
		}
		echoapp_util.DefaultLogger().Info(string(responseData))
	}
	return nil
}

func reportHourTicket() error {
	clientMap := echoapp.ConfigOpts.ReportTicketMap
	for key, options := range clientMap {
		echoapp_util.DefaultLogger().Infof("推送%s流量,com_id:%d", key, options.ComId)

		db, err := app.GetDb("shop")
		if err != nil {
			return errors.Wrap(err, "GetDb")
		}

		now := time.Now()
		toTime := now.Format("2006-01-02 15:04")
		fromTime := now.Format("2006-01-02")
		var count int
		if err := db.Table("tickets").
			Where("com_id = ?", options.ComId).
			Where("used_at > ? and used_at < ?", fromTime, toTime+":59").
			Count(&count).Error; err != nil {
			return errors.Wrap(err, "Query")
		}
		reportDataRequest := &RealPeopleNumber{
			ScenicCode: options.ScenicCode,
			UpTime:     toTime,
			Total:      count,
		}
		data, err := json.Marshal(reportDataRequest)
		echoapp_util.DefaultLogger().Infof("请求数据: %s", string(data))
		if err != nil {
			return errors.Wrap(err, "json.Marshal")
		}

		//入园
		url := options.BaseUrl + "/upload-data/tourist/real-people-number"
		responseData, err := doReportHttpRequest(url, options.AppKey, data)
		if err != nil {
			return errors.Wrap(err, "doReportHttpRequest")
		}

		reportDataResponse := &ReportDataResponse{}
		if err = json.Unmarshal(responseData, reportDataResponse); err != nil {
			return errors.Wrap(err, "doReportHttpRequest")
		}
		echoapp_util.DefaultLogger().Infof("返回结果:%s", string(responseData))

		//出园
		exitUrl := options.BaseUrl + "/upload-data/tourist/real-exit-people-number"
		exitResponseData, err := doReportHttpRequest(exitUrl, options.AppKey, data)
		if err != nil {
			return errors.Wrap(err, "doReportHttpRequest")
		}

		reportExitDataResponse := &ReportDataResponse{}
		if err = json.Unmarshal(exitResponseData, reportExitDataResponse); err != nil {
			return errors.Wrap(err, "doReportHttpRequest")
		}
		echoapp_util.DefaultLogger().Infof("返回结果:%s", string(exitResponseData))
	}
	return nil
}

func reportHourTicketV2() error {
	clientMap := echoapp.ConfigOpts.ReportTicketMap
	for key, options := range clientMap {
		echoapp_util.DefaultLogger().Infof("推送%s流量,com_id:%d", key, options.ComId)

		db, err := app.GetDb("shop")
		if err != nil {
			return errors.Wrap(err, "GetDb")
		}

		now := time.Now()
		hourAgo := now.Add(-time.Minute * 15)
		toTime := now.Format("2006-01-02 15:04:05")
		fromTime := hourAgo.Format("2006-01-02 15:04:05")
		var channelHourReports []*ChannelHourReportV2
		if err := db.Debug().Table("tickets").
			Select("staff_id as channel_id ,count(*) as in_num,count(*) as out_num").
			Where("com_id = ?", options.ComId).
			Where("used_at > ? and used_at < ?", fromTime, toTime).
			Group("staff_id").Scan(&channelHourReports).Error; err != nil {
			return errors.Wrap(err, "Query")
		}
		for _, channel := range channelHourReports {
			channel.RecordTime = toTime
		}
		reportDataRequest := &ReportHourV2{
			LoginName: options.LoginName,
			Pwd:       options.Pwd,
			Data:      channelHourReports,
		}
		data, err := json.Marshal(reportDataRequest)
		if err != nil {
			return errors.Wrap(err, "json.Marshal")
		}
		body := "jsonParams=" + string(data)
		//入园
		url := options.BaseUrl + "/scenic/report"
		echoapp_util.DefaultLogger().Infof("请求数据: %s", string(body))
		responseData, err := doReportHttpRequest(url, options.AppKey, []byte(body))
		if err != nil {
			return errors.Wrap(err, "doReportHttpRequestV2")
		}

		reportDataResponse := &[]ReportDataResponseV2{}
		if err = json.Unmarshal(responseData, reportDataResponse); err != nil {
			return errors.Wrap(err, "doReportHttp response V2:"+string(responseData))
		}
		echoapp_util.DefaultLogger().Infof("返回结果:%s", string(responseData))
	}
	return nil
}

var reportWay string

/**
go run entry/main.go report-ticket  -m hour
go run entry/main.go report-ticket  -m daily
*/
var reportTicketCmd = &cobra.Command{
	Use:   "report-ticket",
	Short: "上报",
	Long:  `上报使用情况`,
	Run: func(cmd *cobra.Command, args []string) {
		echoapp_util.DefaultLogger().Infof("推送景区人流量 method:" + reportWay)
		switch reportWay {
		case "daily":
			if err := reportDailyTicket(); err != nil {
				echoapp_util.DefaultLogger().Error(err)
			} else {
				echoapp_util.DefaultLogger().Infof("推送成功")
			}
		case "hour":
			if err := reportHourTicketV2(); err != nil {
				echoapp_util.DefaultLogger().Error(err)
			} else {
				echoapp_util.DefaultLogger().Infof("推送完成")
			}
		}

	},
}

func init() {
	reportTicketCmd.PersistentFlags().StringVarP(&reportWay, "method", "m", "daily", "访问方法")
	rootCmd.AddCommand(reportTicketCmd)
}
