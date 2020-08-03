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
	"fmt"
	"time"

	"github.com/gw123/echo-app/app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Statistics struct {
	ID       uint
	Date     string `json:"date"`
	TargetId uint   `json:"target_id"`
	Total    int64  `json:"total"`
	Type     string `json:"type"`
}

func statisticDailyHistory() error {
	echoapp_util.DefaultLogger().Infof("统计当天历史记录")
	now := time.Now()
	Hour := now.Format("2006-01-02 15")
	nextHour := now.Add(time.Hour).Format("2006-01-02 15")
	fmt.Println(Hour, nextHour)
	db, err := app.GetDb("user")
	if err != nil {
		return errors.Wrap(err, "db user")
	}
	//var count int
	result := []*Statistics{}
	if err := db.Table("user_history").
		Select("targer_id,type,date(created_at) as date,sum(amount) as total").
		Where("created_at>? and created_at<?", Hour, nextHour).
		Group("target_id,type").
		Find(&result).Error; err != nil {
		return errors.Wrap(err, "db user_his groupby")
	}
	fmt.Println(result)
	for _, val := range result {
		if err := db.Create(val).Error; err != nil {
			return errors.Wrap(err, "db create")
		}
	}
	return nil
}

var statisticHistoryCmd = &cobra.Command{
	Use:   "statistic-history",
	Short: "统计",
	Long:  "按小时统计点击量",
	Run: func(cmd *cobra.Command, args []string) {
		echoapp_util.DefaultLogger().Infof("统计目标点击率")
		statisticDailyHistory()
	},
}

func init() {
	rootCmd.AddCommand(statisticHistoryCmd)
}
