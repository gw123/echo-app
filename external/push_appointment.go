package external

import (
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
"notice": "当日景区游客已达最大承载量，已关闭景区预约通道，暂停对未预订景区门票
游客的接待",
"realtime_tourists": 2000,
"total_tourist": 10000
}
]
响
*/
var appointmentPushUrl = ""

type BookItem struct {
	Label     string `json:"label"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	BookNum   int    `json:"book_num"`
	RemainNum int    `json:"remain_num"`
}

type PushAppointmentRequest struct {
	ScenicID string
	BookList []map[string]BookItem
}

type PushAppointmentResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

// 推送预约信息到旅游局
func DoPushAppointmentRequest(request *PushAppointmentRequest) (*PushAppointmentResponse, error) {
	var response PushAppointmentResponse
	err := DoPost(appointmentPushUrl, request, response)
	if err != nil {
		return nil, errors.Wrap(err, "DoPushAppointmentRequest marshal")
	}

	if response.Code != 0 {
		return nil, errors.Wrap(err, "DoPushAppointmentRequest errMsg "+response.Msg)
	}
	return &response, nil
}
