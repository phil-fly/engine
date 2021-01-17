package email

import (
	"crypto/tls"
	"github.com/go-gomail/gomail"
	"reflect"
)

const ParamTable = ` <tr>
                    <td>
                        {name}：
                        {value}
                    </td>
                </tr>
`



type MailServerSetting struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port int `json:"port"`
	Tls  bool   `json:"tls"`
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