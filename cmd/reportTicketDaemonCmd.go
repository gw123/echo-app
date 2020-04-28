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
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

type ReportDataRequest struct {
}
type ReportDataResponse struct {
}

func reportTicket() error {
	//start := time.Now().Format("2006-01-02 15:04:05")
	clientMap := echoapp.ConfigOpts.ReportTicketMap

	for key, options := range clientMap {
		echoapp_util.DefaultLogger().Infof("推送%s流量,com_id:%d", key, options.ComId)
		reportDataRequest := &ReportDataRequest{}
		data, err := json.Marshal(reportDataRequest)

		if err != nil {
			return errors.Wrap(err, "json.Marshal")
		}

		responseData, err := doReportHttpRequest(options.Addr, data)
		if err != nil {
			return errors.Wrap(err, "doReportHttpRequest")
		}

		reportDataResponse := &ReportDataResponse{}
		if err = json.Unmarshal(responseData, reportDataResponse); err != nil {
			return errors.Wrap(err, "doReportHttpRequest")
		}
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
		startSmsDaemon()
	},
}

func init() {
	rootCmd.AddCommand(reportTicketCmd)
}
