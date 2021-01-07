package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-gomail/gomail"
	"reflect"
	"strings"
)

const OnlineNotice_tpl = `<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
    <style>
        .tag {
            padding: 2px 10px;
            font-size: 12px;
            border-radius: 5px;
            border: 1px solid #fff;
        }

        .danger {
            color: red;
            border-color: red;
        }
        table td{
            word-break:break-all;
        }
    </style>
</head>

<body>
    <div style="max-width: 600px;margin: 0 auto;">
        <div class="head" style="    margin: 0px;
        font-size: 0px;
        vertical-align: top;
        background-color: #1c272b;
        border-bottom: 1px solid #fff;
       padding: 16px;
        text-align: center;">
            <p style="font-size: 14px;color: #fff;font-weight: bold;">{messageType}</p>
        </div>
        <div class="body"
            style="line-height: 24px;background-color: #f8f8f8;color: #424651;font-size: 14px;padding: 20px;">
            <table style="width: 100%;">
                <tr>
                    <td>
                        <h4>
                            <span class="tag danger">高危</span>
                            {title}
                            <span class="system">windows</span>
                        </h4>
                    </td>
                </tr>
                {tplParam}
                 <tr>
                    <td>
                        <p class="msg">
                            <div style="padding:10px 10px 0;border-top:1px solid #ccc;color:#747474;margin-bottom:20px;line-height:1.3em;font-size:12px;">
                    <p>此为系统邮件，请勿回复<br>
                        请保管好您的邮箱，避免账号被他人盗用
                    </p>
                    <p>{Organization}团队</p>
                </div>
                        </p>
                    </td>
                </tr>
            </table>
        </div>
    </div>

</body>
`

const ParamTable = ` <tr>
                    <td>
                        {name}:
                    </td>
					<td>
                        {value}
                    </td>
                </tr>
`




type OnlineNoticeTpl struct {
	serverSetting MailServerSetting `json:"MailServerSetting"`
	doc   DocItem
}

// 文档对象
type DocItem struct {
	MessageType	string	`json:"messageType"`
	Title  string  `json:"title"`  // 通知主题
	Fields []Field `json:"fields"` // 字段列表
	Organization	string	`json:"organization"` // 团队名称
}

// 字段
type Field struct {
	Name  string `json:"name"`  // 字段名称
	Value string `json:"value"` // 字段值
}

func (self *OnlineNoticeTpl) SetTitle(title string) {
	self.doc.Title = title
}

func (self *OnlineNoticeTpl) SetOrganization(organization string) {
	self.doc.Organization = organization
}

func (self *OnlineNoticeTpl) SetMessageType(messageType string) {
	self.doc.MessageType = messageType
}

func (self *OnlineNoticeTpl) SetContent(param interface{}) {
	self.doc.Fields = createFields(param)
}

func (self *OnlineNoticeTpl) MailServerSetting(ServerSetting MailServerSetting) {
	self.serverSetting = ServerSetting
}

func (m *OnlineNoticeTpl) RenderParam(v Field) string {
	ts := ParamTable

	ts = strings.Replace(ts, "{name}", v.Name, 1)
	ts = strings.Replace(ts, "{value}", v.Value, 1)
	return ts
}

func (m *OnlineNoticeTpl) RenderContent(v DocItem) string {
	ts := OnlineNotice_tpl
	var tplParams string = ""
	ts = strings.Replace(ts, "{messageType}", v.MessageType, 1)
	ts = strings.Replace(ts, "{title}", v.Title, 1)
	if len(v.Fields) > 0 {
		for _, item := range v.Fields {
			tpl := m.RenderParam(item)
			tplParams = fmt.Sprintf("%s%s", tplParams, tpl)
		}
	}
	ts = strings.Replace(ts, "{tplParam}", tplParams, 1)
	ts = strings.Replace(ts, "{Organization}", v.Organization, 1)
	return ts
}

func (self *OnlineNoticeTpl)SendMail(toUser string) error {
	emailBody:= self.RenderContent(self.doc)
	if emailBody == ""{
		return errors.New("html is nil")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", self.serverSetting.User)
	m.SetHeader("To", toUser)
	m.SetHeader("Subject", "攻击者上线通知") // 邮件标题
	m.SetBody("text/html", emailBody) // 邮件内容
	// m.Attach("/home/Alex/lolcat.jpg") //附件

	d := gomail.NewDialer(self.serverSetting.Host, self.serverSetting.Port, self.serverSetting.User, self.serverSetting.Pass)
	if self.serverSetting.Tls {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func createFields(param interface{}) []Field {
	if param == nil {
		return nil
	}
	fields := make([]Field, 0)
	val := reflect.ValueOf(param)
	if !val.IsValid() {
		panic("not valid")
	}
	//fmt.Println(val.Kind())
	if val.Kind() == reflect.Slice {
		if val.Len() > 0 {
			val = val.Index(0)
		} else {
			return nil
		}
	}
	for val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}
	typ := val.Type()
	cnt := val.NumField()

	for i := 0; i < cnt; i++ {
		ty := typ.Field(i)
		field := Field{
			Name:  ty.Tag.Get("doc"),
			Value: val.Field(i).Interface().(string),
		}
		fields = append(fields, field)
	}
	return fields
}
