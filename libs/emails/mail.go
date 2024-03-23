package emails

import (
	"github.com/pkg/errors"

	"github.com/go-gomail/gomail"
)

type MailService interface {
	SendHtmlMail(to []string, subject string, body string) error
	SendMail(to []string, subject, body, filename string) error
}

type GoMail struct {
	host     string
	port     int
	username string
	password string
	dialer   *gomail.Dialer
}

func NewGoMail(host string, port int, username, password string) *GoMail {
	dialer := gomail.NewDialer(host, port, username, password)
	return &GoMail{
		host:     host,
		port:     port,
		username: username,
		password: password,
		dialer:   dialer,
	}
}

func (g GoMail) SendHtmlMail(to []string, subject string, body string) error {
	m := gomail.NewMessage()
	// 收件人可以有多个，故用此方式
	m.SetHeader("To", to...)
	// 发件人
	// 第三个参数为发件人别名，如"李大锤"，可以为空（此时则为邮箱名称）
	m.SetAddressHeader("From", g.username, "")
	// 主题
	m.SetHeader("Subject", subject)
	// 正文
	m.SetBody("text/html", body)

	if err := g.dialer.DialAndSend(m); err != nil {
		return errors.Wrap(err, "dialAndSend")
	}
	return nil
}

func (g GoMail) SendMail(to []string, subject, body string, filename string) error {
	m := gomail.NewMessage()
	// 收件人可以有多个，故用此方式
	m.SetHeader("To", to...)
	// 发件人
	// 第三个参数为发件人别名，如"李大锤"，可以为空（此时则为邮箱名称）
	m.SetAddressHeader("From", g.username, "")
	// 主题
	m.SetHeader("Subject", subject)
	// 正文
	m.SetBody("text/html", body)
	m.Attach(filename)
	if err := g.dialer.DialAndSend(m); err != nil {
		return errors.Wrap(err, "dialAndSend")
	}
	return nil
}
