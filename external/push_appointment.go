package external

import (
	"context"

	"github.com/gw123/glog"

	"github.com/pkg/errors"
)

/***
[
	{
	"scenic_id": "130",
	"book_list": [
		{
			"2020-05-28": [
				{
				"label": "08:00-09:00",
				"start_time": 1590278400,
				"end_time": 1590282000,
				"book_num": 8,
				"remain_num": 22
				}
			]
		}
	],
	"scenic_status": 1,
	"notice": "当日景区游客已达最大承载量，已关闭景区预约通道，暂停对未预订景区门票游客的接待",
	"realtime_tourists": 2000,
	"total_tourist": 10000
	}
]
*/
var appointmentPushUrl = "http://bigd.oindata.cn/data-api/api/bookData"

type BookItem struct {
	Label     string `json:"label"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
	BookNum   int    `json:"book_num"`
	RemainNum int    `json:"remain_num"`
}

type PushAppointmentRequest struct {
	ScenicID         string                  `json:"scenic_id"`
	BookList         []map[string][]BookItem `json:"book_list"`
	ScenicStatus     int                     `json:"scenic_status"`
	Notice           string                  `json:"notice"`
	RealtimeTourists int                     `json:"realtime_tourists"`
	TotalTourist     int                     `json:"total_tourist"`
}

//type PushAppointmentRequest = []PushAppointmentReques
type PushAppointmentResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// 推送预约信息到旅游局
func DoPushAppointmentRequest(ctx context.Context, request []*PushAppointmentRequest) (*PushAppointmentResponse, error) {
	glog.ExtractEntry(ctx).WithField("request", request).Infof("DoPushAppointmentRequest request")
	var response PushAppointmentResponse
	err := DoPost(appointmentPushUrl, request, &response)
	if err != nil {
		glog.ExtractEntry(ctx).WithError(err).Infof("DoPushAppointmentRequest err")
		return nil, errors.Wrap(err, "DoPushAppointmentRequest marshal")
	}
	glog.ExtractEntry(ctx).WithField("response", response).Infof("DoPushAppointmentRequest response")

	if response.Code != 0 {
		return nil, errors.Wrap(err, "DoPushAppointmentRequest errMsg "+response.Msg)
	}
	return &response, nil
}
