package echoapp_util

import (
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

// type Glog struct {
// 	//glog *glog.
// 	clog *log.Logger
// }

// func (l *Glog) Info(msg string, keysAndValues ...interface{}) {
// 	l.clog.WithFields(log.Fields{
// 		"data": keysAndValues,
// 	}).Info(msg)
// }
// func (l *Glog) Error(err error, msg string, keysAndValues ...interface{}) {
// 	l.clog.WithFields(log.Fields{
// 		"msg":  msg,
// 		"data": keysAndValues,
// 	}).Warn(msg)
// }
// func NewGlog() cron.Logger {
// 	return &Glog{}
// }

// //TestCorn tset
// func testCorn() {
// 	//c := cron.New()
// 	glog := NewGlog()
// 	//// 可以配置如果当前任务正在进行，那么跳过
// 	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(glog)))
// 	i := 1
// 	c.AddFunc("*/1 * * * * ", func() {
// 		fmt.Println("每分钟执行一次", i)
// 		i++
// 	})
// 	//每半个小时
// 	c.AddFunc("30 * * * *", func() { fmt.Println("Every hour on the half hour") })
// 	// 在凌晨3点到6点，晚上8点到11点之间
// 	c.AddFunc("30 3-6,20-23 * * *", func() { fmt.Println(".. in the range 3-6am, 8-11pm") })
// 	//每天东京时间4:30
// 	c.AddFunc("CRON_TZ=Asia/Tokyo 30 04 * * *", func() { fmt.Println("Runs at 04:30 Tokyo time every day") })
// 	// 每小时一次，从现在开始一小时
// 	c.AddFunc("@hourly", func() { fmt.Println("Every hour, starting an hour from now") })
// 	// 每小时三十分，从现在开始一小时三十分开始
// 	c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty, starting an hour thirty from now") })
// 	c.Start()
// 	// Funcs are invoked in their own goroutine, asynchronously.
// 	// Funcs may also be added to a running Cron
// 	c.AddFunc("@daily", func() { fmt.Println("Every day") })
// 	// Inspect the cron job entries' next and previous run times.
// 	//inspect(c.Entries())

// 	//c.Stop() // Stop the scheduler (does not stop any jobs already running).
// 	time.Sleep(time.Hour * 4)

// }
//ParseCronString 解析cron字符串，仅限 标准格式 “* * * * *”,判断start 是否为 cron 任务的下一个开始时间
var specParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

func ParseCronString(crontab string) (string, error) {
	if len(crontab) == 0 {
		return "", errors.New("nil string")
	}

	sched, err := specParser.Parse(crontab)
	if err != nil {
		return "", errors.Wrap(err, "specParser.Parse:"+crontab)
	}
	nowTime := time.Now()
	nTimeZero := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, nowTime.Location())
	nextTimeZero := nTimeZero.AddDate(0, 0, 1)
	nextCronTime := sched.Next(nTimeZero)
	if nextCronTime == nextTimeZero {
		return nTimeZero.Format(TIME_LAYOUT), nil
	} else if nextCronTime.Before(nextTimeZero) {
		return nextCronTime.Format(TIME_LAYOUT), nil
	}
	return "", errors.New("parse failed")

}
func ParseCronStringByStartTime(crontab string, start, end, queryTime time.Time) (bool, error) {
	if len(crontab) == 0 {
		return false, errors.New("nil string")
	}
	specParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sched, err := specParser.Parse(crontab)
	if err != nil {
		return false, errors.Wrap(err, crontab)
	}
	nextCronTime := sched.Next(queryTime)

	if nextCronTime.Before(start) || nextCronTime.After(end) {
		return false, errors.New("before")
	}
	return true, nil

}

const (
	TIME_LAYOUT = "2006-01-02 15:04"
)

func ParseWithLocation(locationName string, timeStr string) (time.Time, error) {
	location, err := time.LoadLocation(locationName)
	if err != nil {
		println(err.Error())
		return time.Time{}, errors.Wrap(err, "time.LoadLocation")
	}

	lt, err := time.ParseInLocation(TIME_LAYOUT, timeStr, location)
	//fmt.Println(location, lt)
	return lt, nil

}
