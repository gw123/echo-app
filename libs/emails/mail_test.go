package emails

import (
	"testing"
)

var (
	serverHost = "smtp.ym.163.com"
	serverPort = 25
	fromEmail  = "robot@xytschool.com"
	fromPasswd = "oWVeEpdjoc"
)

func TestGoMail_SendHtmlMail(t *testing.T) {
	mail := NewGoMail(serverHost, serverPort, fromEmail, fromPasswd)
	subject := "这是主题3"
	body := `这是正文<br>
            <h3>这是标题</h3>
             Hello <a href = "http://www.latelee.org">主页</a><br>`
	err := mail.SendHtmlMail(fromEmail, []string{"963353840@qq.com"}, subject, body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send success")
}

func TestGoMail_SendMail(t *testing.T) {
	mail := NewGoMail(serverHost, serverPort, fromEmail, fromPasswd)
	subject := "这是主题22"
	body := `这是正文<br>
            <h3>这是标题</h3>
             Hello <a href = "http://www.latelee.org">主页</a><br>`
	filename := "/tmp/test.txt"
	err := mail.SendMail(fromEmail, []string{"963353840@qq.com"}, subject, body, filename)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("send success")
}
