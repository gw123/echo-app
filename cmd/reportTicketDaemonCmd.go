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
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Data struct {
	InNum      string `json:"inNum"`
	OutNum     string `json:"outNum"`
	ChannelId  string `json:"channelId"`
	RecordTime string `json:"recordTime"`
}
type ReportDataRequest struct {
	LoginName string `json:"loginName"`
	Pwd       string `json:"pwd"`
	Data      []Data `json:"data"`
}
type ReportDataResponse struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func reportTicket() error {
	//start := time.Now().Format("2006-01-02 15:04:05")
	clientMap := echoapp.ConfigOpts.ReportTicketMap
	for key, options := range clientMap {
		echoapp_util.DefaultLogger().Infof("推送%s流量,com_id:%d", key, options.ComId)

		var datas = []Data{
			{
				InNum:      "1111",
				OutNum:     "999",
				ChannelId:  "01",
				RecordTime: time.Now().Format("2006-01-02 15:04:05"),
			},
			{
				InNum:      "2222",
				OutNum:     "1999",
				ChannelId:  "01",
				RecordTime: time.Now().Add(time.Hour).Format("2006-01-02 15:04:05"),
			},
		}
		reportDataRequest := &ReportDataRequest{
			LoginName: options.Username,
			Pwd:       options.Password,
			Data:      datas,
		}

		/*reportDataRequest := &ReportDataRequest{}
		var ctx echo.Context
		if err := ctx.Bind(reportDataRequest); err != nil {
			panic(err)
		}*/
		data, err := json.Marshal(reportDataRequest)
		//fmt.Println(data)
		if err != nil {
			return errors.Wrap(err, "json.Marshal")
		}

		responseData, err := doReportHttpRequest(options.Addr, data)
		//fmt.Println(responseData)
		if err != nil {
			return errors.Wrap(err, "doReportHttpRequest")
		}

		reportDataResponse := &ReportDataResponse{}
		if err = json.Unmarshal(responseData, reportDataResponse); err != nil {
			return errors.Wrap(err, "doReportHttpRequest")
		}
		fmt.Println(reportDataResponse.Msg)
	}
	return nil
}

func doReportHttpRequest(url string, data []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->http.NewRequest")
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->Do")
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	//fmt.Println(data)
	if err != nil {
		return nil, errors.Wrap(err, "DoRequest->ReadAll")
	}
	if res.StatusCode == http.StatusOK {
		return data, errors.Wrapf(err, "StatusCode:%d", res.StatusCode)
	}
	return data, nil
}

// serverCmd represents the server command
var reportTicketCmd = &cobra.Command{
	Use:   "report-ticket",
	Short: "上报",
	Long:  `上报使用情况`,
	Run: func(cmd *cobra.Command, args []string) {
		echoapp_util.DefaultLogger().Infof("推送景区人流量")
		reportTicket()
	},
}

func init() {
	rootCmd.AddCommand(reportTicketCmd)
}
