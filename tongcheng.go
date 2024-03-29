package echoapp

// TongchengRequestHead
type TongchengRequestHead struct {
	Sign      string `json:"sign"`
	UserId    string `json:"user_id"`
	Method    string `json:"method"`
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

// TongchengRequest 基础请求结构
type TongchengRequest struct {
	RequestHead        TongchengRequestHead `json:"requestHead"`
	RawRequestBody     string               `json:"-"`
	EncryptRequestBody string               `json:"requestBody"`
}

// TongchengResponse 响应结构
type TongchengResponse struct {
	ResponseHead struct {
		ResCode   string `json:"res_code"`
		ResMsg    string `json:"res_msg"`
		Timestamp int    `json:"timestamp"`
	} `json:"responseHead"`
	ResponseBody string `json:"responseBody"`
}

//  CheckTicketRequestBody
type CheckTicketRequestBody struct {
	Tickets        uint   `json:"tickets"`
	ComID          uint   `json:"com_id"`
	OrderSerialId  string `json:"orderSerialId"`
	PartnerOrderId string `json:"partnerOrderId"`
	ConsumeDate    string `json:"consumeDate"`
}

type SyncPartnerCodeRequestBody struct {
	Tickets        uint   `json:"tickets"`
	OrderSerialId  string `json:"orderSerialId"`
	PartnerOrderId string `json:"partnerOrderId"`
	PartnerCode    string `json:"partnerCode"`
}

type CheckTicketJob struct {
	BaseMqMsg
	CheckTicketRequestBody
}

type SyncPartnerCodeJob struct {
	BaseMqMsg
	SyncPartnerCodeRequestBody
}

type TongchengService interface {
	//核销门票通知第三方
	SyncPartnerCode(comID uint, info *SyncPartnerCodeRequestBody) error
	CheckTicket(comID uint, info *CheckTicketRequestBody) error
}
