package echoapp

type SendMessageOptions struct {
	ComId         uint     `json:"com_id"`
	PhoneNumbers  []string `json:"phone_numbers"`
	Type          string   `json:"type"`
	TemplateParam string   `json:"template_param"`
}

type SmsChannel struct {
	Type         string `json:"type"`
	Channel      string `json:"channel"`
	Key          string `json:"key"`
	Secret       string `json:"secret"`
	SignName     string `json:"sign_name"`
	TemplateCode string `json:"template_code"`
}

type SendMessageJob struct {
	BaseMqMsg
	SendMessageOptions
}

type SmsService interface {
	CheckVerifyCode(comId uint,phone string, code string) bool
	SendVerifyCodeSms(comId uint, phone string ,code string) error
	SendMessage(opt *SendMessageOptions) error
}
