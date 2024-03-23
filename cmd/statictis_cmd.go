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

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/gw123/echo-app/app/app_components"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	staHour = 1
	//comId =14
)

type Statistics struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Date      string `json:"date"`
	TargetId  int    `json:"target_id"`
	Total     int64  `json:"total"`
	Type      string `json:"type"`
}

func (*Statistics) TableName() string {
	return "Statistic1s"
}
func statisticDailyHistory(start, end string) error {
	echoapp_util.DefaultLogger().Infof("统计当天历史记录")
	// now := time.Now()
	// nowStr := now.Format("2006-01-02 15:04:05")
	// nowHourStr := now.Format("2006-01-02 15")
	// nowDaystr := now.Format("2006-01-02")
	// // start := now.Add(time.Hour * (staHour))
	// // startStr := start.Format("2006-01-02 15:04:05")
	// // startHourStr:=start.Format("2006-01-02 15")
	// // startDayStr:=start.Format("2006-01-02")
	db, err := app.GetDb("user")
	if err != nil {
		return errors.Wrap(err, "db user")
	}
	if start > end {
		start, end = end, start
	}
	result := []*Statistics{}
	err = db.Table("user_history").
		Where("created_at>=? and created_at<=? ", start, end). //统计当天0时到现在
		Select("count(*) as total,target_id,type").
		Group("target_id,type").
		Find(&result).Error
	if err != nil {
		return errors.Wrap(err, "db user_his groupby")
	}
	db.AutoMigrate(&Statistics{})
	for _, val := range result {
		val.Date = fmt.Sprintf("%s--%s", start, end)
		//fmt.Println(val)
		if err := db.Save(val).Error; err != nil {
			return errors.Wrap(err, "db create")
		}

	}
	return nil
}
func statisticCompanySales(start, end string) error {
	echoapp_util.DefaultLogger().Infof("统计Company销量")
	salesStatistic := []*echoapp.CompanySalesSatistic{}
	db, err := app_components.GetShopDb()
	if err != nil {
		return errors.Wrap(err, ".GetDb")
	}
	if err := db.Table("orders").
		Where("created_at>=? and created_at<=?", start, end).
		Select("com_id,sum(total) as all_sales_total").Group("com_id").
		Find(&salesStatistic).Error; err != nil {
		return errors.Wrap(err, "query")
	}
	fmt.Println(salesStatistic)
	db.AutoMigrate(&echoapp.CompanySalesSatistic{})
	for _, val := range salesStatistic {
		val.Date = fmt.Sprintf("%s--%s", start, end)
		if err := db.Create(val).Error; err != nil {
			return errors.Wrap(err, "Create")
		}
	}

	return nil
}
func statisticComGoodsSales(start, end string, comId uint) error {
	echoapp_util.DefaultLogger().Infof("统计商品销量")
	salesStatistic := []*echoapp.GoodsSalesSatistic{}
	db, err := app_components.GetShopDb()
	if err != nil {
		return errors.Wrap(err, ".GetDb")
	}
	if err := db.Table("orders").Where("com_id=?", comId).
		Where("created_at>=? and created_at<=?", start, end).
		Select("goods_id,sum(total) as goods_sales").Group("goods_id").
		Find(&salesStatistic).Error; err != nil {
		return errors.Wrap(err, "query")
	}
	db.AutoMigrate(&echoapp.GoodsSalesSatistic{})
	for _, val := range salesStatistic {
		val.Date = fmt.Sprintf("%s--%s", start, end)
		val.ComId = comId
		if err := db.Create(val).Error; err != nil {
			return errors.Wrap(err, "Create")
		}
	}
	return nil
}

var statisticsWay string
var statisticHistoryCmd = &cobra.Command{
	Use:   "statistic",
	Short: "统计",
	Long:  "按小时统计",
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		nowStr := now.Format("2006-01-02 15:04:05")
		//nowHourStr := now.Format("2006-01-02 15")
		nowDaystr := now.Format("2006-01-02")
		// start := now.Add(time.Hour * (staHour))
		// startStr := start.Format("2006-01-02 15:04:05")
		// startHourStr:=start.Format("2006-01-02 15")
		// startDayStr:=start.Format("2006-01-02")
		testDate := "2020-04-27"
		echoapp_util.DefaultLogger().Infof("统计时间：%s--%s", nowDaystr, nowStr)
		switch statisticsWay {
		case "history":
			if err := statisticDailyHistory(testDate, nowStr); err != nil {
				echoapp_util.DefaultLogger().Error(err)
			} else {
				echoapp_util.DefaultLogger().Infof("统计完成")
			}
		case "company":
			if err := statisticCompanySales(testDate, nowStr); err != nil {
				echoapp_util.DefaultLogger().Error(err)
			} else {
				echoapp_util.DefaultLogger().Infof("统计完成")
			}
		case "goods":
			var comId uint

			fmt.Printf("请输入要统计的公司ID:")
			fmt.Scanln(&comId)
			if err := statisticComGoodsSales(nowDaystr, nowStr, comId); err != nil {
				echoapp_util.DefaultLogger().Error(err)
			} else {
				echoapp_util.DefaultLogger().Infof("统计完成")
			}
		}
	},
}

func init() {
	statisticHistoryCmd.PersistentFlags().StringVarP(&statisticsWay, "method", "m", "history", "访问方法")
	RootCmd.AddCommand(statisticHistoryCmd)

}
