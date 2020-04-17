package echoapp

import "github.com/labstack/echo"

type SendMessageOptions struct {
	Token         string   `json:"token"`
	PhoneNumbers  []string `json:"phone_numbers"`
	SignName      string   `json:"sign_name"`
	TemplateCode  string   `json:"template_code"`
	TemplateParam string   `json:"template_param"`
}

type SmsService interface {
	SendMessage(ctx echo.Context, opt SendMessageOptions) error
}
