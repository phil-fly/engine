package wechatMessage

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// 文档对象
type DocItem struct {
	Title  string  `json:"title"`  // 企微通知主题
	Fields []Field `json:"fields"` // 字段列表
}

// 字段
type Field struct {
	Name  string    `json:"name"`  // 字段名称
	Color string `json:"color"` // 字段颜色
	Value string      `json:"Value"` // 字段值
}




type Message struct {
	Token		string
	doc         DocItem
}

func (m *Message) SetToken(Token string) {
	m.Token = Token
}

func (m *Message) SetTitle(Title string) {
	m.doc.Title = Title
}

func (m *Message) SetContent(param interface{}) {
	m.doc.Fields = createFields(param)
}


func (m *Message) RenderParam(v Field) string {
	ts := ParamTable

	ts = strings.Replace(ts, "{name}", v.Name, 1)
	ts = strings.Replace(ts, "{color}", v.Color, 1)
	ts = strings.Replace(ts, "{value}", v.Value, 1)
	return ts
}

func (m *Message) RenderContent(v DocItem) string {
	ts := Content
	var tplParams string = ""
	ts = strings.Replace(ts, "{title}", v.Title, 1)
	if len(v.Fields) >0{
		for _, item := range v.Fields {
			tpl := m.RenderParam(item)
			tplParams = fmt.Sprintf("%s%s", tplParams, tpl)
		}
	}
	ts = strings.Replace(ts, "{tplParam}", tplParams, 1)
	return ts
}

func (m *Message) RenderTplPage() string {
	ts := TplPage
	content := m.RenderContent(m.doc)
	ts = strings.Replace(ts, "{content}", content, 1)
	return ts
}

func (m *Message) Send() error {
	if m.Token == "" {
		return errors.New("Token is nil.")
	}
	meg := m.RenderTplPage()
	send_url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + m.Token
	//print(send_url)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", send_url, bytes.NewBuffer([]byte(meg)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//fmt.Println("response Status:", resp.Status)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	return nil
}

type MessageWechat struct {
	DangerLevel string `json:"DangerLevel" doc:"危险等级" color:"warning"`
	HostName string `json:"HostName" doc:"主机名" color:"comment"`
	Os string `json:"Os" doc:"系统类型" color:"comment"`
	ConnectAddr string `json:"ConnectAddr" doc:"连接地址" color:"comment"`
	Tips string `json:"Tips" doc:"提示" color:"comment"`
	Addr string `json:"Addr" doc:"地理位置" color:"comment"`
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
			Name:        ty.Tag.Get("doc"),
			Color:       ty.Tag.Get("color"),
			Value:   	 val.Field(i).Interface().(string),
		}
		fields = append(fields, field)
	}
	return fields
}