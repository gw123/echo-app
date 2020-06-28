package echoapp

type SendMessageOptions struct {
	Token         string   `json:"token"`
	ComId         int      `json:"com_id"`
	PhoneNumbers  []string `json:"phone_numbers"`
	SignName      string   `json:"sign_name"`
	TemplateCode  string   `json:"template_code"`
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
	SendMessage(opt *SendMessageOptions) error
}
