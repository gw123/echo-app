package echoapp_middlewares

import (
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func NewLimitMiddlewares(skipper middleware.Skipper, limitPreSecond int, maxWorker int) echo.MiddlewareFunc {
	if limitPreSecond == 0 {
		limitPreSecond = 100
	}
	if maxWorker == 0 {
		maxWorker = 100
	}

	currentWorkerNum := make(chan uint8, maxWorker)
	for i := 0; i < maxWorker; i++ {
		currentWorkerNum <- 1
	}
	var lastSecondTime int64 = 0
	lastSecondRequestNum := 0
	mutex := sync.Mutex{}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			mutex.Lock()
			if lastSecondTime == time.Now().Unix() {
				lastSecondRequestNum++
				if lastSecondRequestNum > limitPreSecond {
					mutex.Unlock()
					echoapp_util.ExtractEntry(c).Warn("触发服务端秒级限流")
					return c.JSON(http.StatusTooManyRequests, "触发服务端限流")
				}
				mutex.Unlock()
			} else {
				//todo 这里有并发问题,如果直接加锁比较重考虑一种性能高的方式
				lastSecondTime = time.Now().Unix()
				lastSecondRequestNum = 0
				mutex.Unlock()
			}

			//
			select {
			case <-currentWorkerNum:
				defer func() {
					currentWorkerNum <- 1
				}()
			default:
				echoapp_util.ExtractEntry(c).Warn("触发服务端工作协程上限")
				return c.JSON(http.StatusTooManyRequests, "触发服务端限流")
			}
			echoapp_util.AddField(c, "work_num", strconv.Itoa(len(currentWorkerNum)))
			//echoapp_util.ExtractEntry(c).Infof("currentWokerNum: %d", len(currentWorkerNum))
			if skipper(c) {
				return next(c)
			}
			return next(c)
		}
	}
}
