package echoapp

import (
	"github.com/labstack/echo"
)

type CheckTicketJob struct {
	BaseMqMsg
	ComId string `json:"com_id"`
	TongchengRequestBody
}

type TongchengRequest struct {
	RequestHead        TongchengRequestHead `json:"requestHead"`
	RequestBody        TongchengRequestBody `json:"_"`
	RawRequestBody     string               `json:"-"`
	EncryptRequestBody string               `json:"requestBody"`
}

type TongchengRequestHead struct {
	Sign      string `json:"sign"`
	UserId    string `json:"user_id"`
	Method    string `json:"method"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

type TongchengRequestBody struct {
	Tickets        int    `json:"tickets"`
	OrderSerialId  string `json:"orderSerialId"`
	PartnerOrderId string `json:"partnerOrderId"`
	ConsumeDate    string `json:"consumeDate"`
}

type TongchengConsumeNoticeRequest struct {
	RequestHead TongchengRequestHead `json:"requestHead"`
	RequestBody TongchengRequestBody `json:"requestBody"`
}

type TongchengResponse struct {
	ResponseHead struct {
		ResCode   string `json:"res_code"`
		ResMsg    string `json:"res_msg"`
		Timestamp int    `json:"timestamp"`
	} `json:"responseHead"`
	ResponseBody string `json:"responseBody"`
}

type TongchengService interface {
	//核销门票通知第三方
	CheckTicket(ctx echo.Context, info CheckTicketJob) error
}
