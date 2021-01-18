package external

import (
	"github.com/pkg/errors"
)

var push_checked_ticket_hour = "http://qmt.lyta.com.cn/upload-data/tourist/real-people-number"
var push_checked_ticket_daily = "http://qmt.lyta.com.cn/upload-data/tourist/real-gate-day"

type ChannelReport struct {
	InNum      int    `json:"inNum" gorm:"column:in_num" `
	OutNum     int    `json:"outNum" gorm:"column:out_num"`
	ChannelId  string `json:"channelId" gorm:"column:channel_id"`
	RecordTime string `json:"recordTime"`
}

type PushCheckTicketRequest struct {
	LoginName string           `json:"loginName"`
	Pwd       string           `json:"pwd"`
	Data      []*ChannelReport `json:"data"`
}

type PushCheckTicketResponse struct {
	Type int    `json:"type"`
	Msg  string `json:"msg"`
}

// 推送核销门票消息
func DoPushTicketRequest(request *PushCheckTicketRequest) (*PushCheckTicketResponse, error) {
	var response PushCheckTicketResponse
	err := DoPost(push_checked_ticket_hour, request, response)
	if err != nil {
		return nil, errors.Wrap(err, "DoPushTicketRequest marshal")
	}

	if response.Type != 0 {
		return nil, errors.Wrap(err, "DoPushTicketRequest errMsg "+response.Msg)
	}
	return &response, nil
}
