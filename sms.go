package echoapp

type SendMessageOptions struct {
	Token         string
	PhoneNumbers  []string
	SignName      string
	TemplateCode  string
	TemplateParam string
}

type SmsService interface {
	SendMessage(opt SendMessageOptions) error
}
