package crontabs

import (
	"context"
	"os"
	"os/signal"
	"time"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/gw123/echo-app/external"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

// type Book struct {
// 	Label     string    `json:"label"`
// 	StartTime time.Time `json:"start_time"`
// 	EndTime   time.Time `json:"end_time"`
// 	BookNum   int64     `json:"book_num"`
// 	RemainNum int64     `json:"remain_num"`
// }
// type BookDataRequest struct {
// 	SenicID  string `json:"senic_id"`
// 	BookList []Book `json:"book_list"`

// 	ScenicStatus     int    `json:"scenic_status"`
// 	Notice           string `json:"notice"`
// 	RealtimeTourists int    `json:"realtime_tourists"`
// 	TotalTourist     int    `json:"total_tourist"`
// }

// type BookDataInitRequest struct {
// 	SenicID     string                 `json:"senic_id"`
// 	BooktimeSet []*echoapp.Appointment `json:"booktime_set"`
// 	Label       string                 `json:"label"`
// 	StartClock  string                 `json:"startClock"`
// 	EndTClock   string                 `json:"endClock"`
// 	MaxBook     int64                  `json:"maxBook"`
// 	MaxCapacity int64                  `json:"max_capacity"`
// }

type Job struct {
	Shut chan int
}

func (j *Job) Run() {
	reportBookingPassengerFlow()
}
func cronJobs() {
	c := cron.New()
	//c.AddFunc("*/15 8-16 * * *", reportBookingPassengerFlow)
	job := Job{make(chan int, 1)}
	c.AddJob("*/15 8-16 * * *", &job)
	c.Start()
	defer c.Stop()
	select {
	case <-job.Shut:
		return
	}

}

var (
	TimeLayout       = "2006-01-02"
	TimeLayoutL      = "2006-01-02 15:04"
	Max              = 100
	StatusCode       = 1
	Notice           = "Notice:"
	RealtimeTourists = 2222
	TotalTourist     = 1232323
)

func reportBookingPassengerFlow() {
	clientMap := echoapp.ConfigOpts.BookingData
	pushARequestArr := []*external.PushAppointmentRequest{}
	for key, options := range clientMap {
		echoapp_util.DefaultLogger().Infof("推送%s预约量,com_id:%d", key, options.ComId)

		db, err := app.GetDb("shop")
		if err != nil {
			echoapp_util.DefaultLogger().Error(err)
			return
		}

		now := time.Now()
		hourAgo := now.Add(-time.Minute * 15)
		toTime := now.Format(TimeLayoutL)
		fromTime := hourAgo.Format(TimeLayoutL)
		var bookNum int
		if err := db.Debug().Table("appointments").
			Where("com_id = ?", options.ComId).
			Where("created_at between ? and ?", fromTime, toTime).
			Count(&bookNum).Error; err != nil {
			echoapp_util.DefaultLogger().Error(err)
			return
		}
		label := fromTime[11:] + "-" + toTime[11:]
		//label := strconv.Itoa(hourAgo.Hour()) + ":" + strconv.Itoa(hourAgo.Minute()) + "-" + strconv.Itoa(now.Hour()) + ":" + strconv.Itoa(now.Minute())
		bookItem := external.BookItem{
			Label:     label,
			StartTime: int(hourAgo.Unix()),
			EndTime:   int(now.Unix()),
			BookNum:   bookNum,
			//RemainNum: Max - bookNum,
		}
		bookMap := make(map[string][]external.BookItem)
		bookMap[now.Format(TimeLayout)] = append(bookMap[now.Format(TimeLayout)], bookItem)
		bookList := []map[string][]external.BookItem{}
		bookList = append(bookList, bookMap)
		pushAppointmentRequest := &external.PushAppointmentRequest{
			ScenicID:         options.ScenicCode,
			BookList:         bookList,
			ScenicStatus:     StatusCode,
			Notice:           Notice,
			RealtimeTourists: RealtimeTourists,
			TotalTourist:     TotalTourist,
		}

		pushARequestArr = append(pushARequestArr, pushAppointmentRequest)

	}
	res, err := external.DoPushAppointmentRequest(context.Background(), pushARequestArr)
	if err != nil {
		echoapp_util.DefaultLogger().Error(err)
		return
	}
	echoapp_util.DefaultLogger().Infof("返回结果:%v", res)

}

var AppointmentCmd = &cobra.Command{
	Use:   "appointment",
	Short: "报告预订客流",
	Long:  `定时任务，每15分上报预订客流`,
	Run: func(cmd *cobra.Command, args []string) {
		echoapp_util.DefaultLogger().Info("推送景区预订客流")
		quit := make(chan os.Signal, 1)
		cronJobs()
		signal.Notify(quit, os.Interrupt)
		<-quit
	},
}
