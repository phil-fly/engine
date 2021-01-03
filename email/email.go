package email

import (
	"crypto/tls"
	"github.com/go-gomail/gomail"
)

type MailServerSetting struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port int `json:"port"`
	Tls  bool   `json:"tls"`
}

func (self *MailServerSetting)sendMail(){
	m := gomail.NewMessage()
	m.SetHeader("From", self.User)
	m.SetHeader("To", "XXXOOO@163.com","XXXOOO@qq.com")
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan") //抄送
	m.SetHeader("Subject", "测试") // 邮件标题
	m.SetBody("text/html", "this is 测试") // 邮件内容
	// m.Attach("/home/Alex/lolcat.jpg") //附件

	d := gomail.NewDialer(self.Host, self.Port, self.User, self.Pass)
	if self.Tls {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}