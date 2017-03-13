package dinghook

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

const (
	// DingAPIURL api 地址
	DingAPIURL = `https://oapi.dingtalk.com/robot/send?access_token=`
)

// Result 发送结果
// Success true 成功，否则失败
// ErrMsg 错误信息，如果是钉钉接口错误，会返回钉钉的错误信息，否则返回内部err信息
// ErrCode 钉钉返回的错误码
type Result struct {
	Success bool
	// ErrMsg 错误信息
	ErrMsg string `json:"errmsg"`
	// 错误码
	ErrCode int `json:"errcode"`
}

// Ding 钉钉消息发送实体
type Ding struct {
	AccessToken string // token
}

// Send 发送消息
func (ding Ding) Send(message interface{}) Result {
	if ding.AccessToken == "" {
		return Result{ErrMsg: "access token is required"}
	}

	switch message.(type) {
	case *Message:
	case Message:
	case Link:
	case *Link:
	case Markdown:
	case *Markdown:
	default:
		return Result{ErrMsg: "not support message type"}
	}

	var err error

	// 检查必填项目
	if err = validator.New().Struct(message); err != nil {
		return Result{ErrMsg: "field valid errror: " + err.Error()}
	}

	var paramsMap map[string]interface{}
	if m, ok := message.(Message); ok {
		paramsMap = convertMessage(m)
	} else if m, ok := message.(*Message); ok {
		paramsMap = convertMessage(*m)
	} else if m, ok := message.(Link); ok {
		paramsMap = convertLink(m)
	} else if m, ok := message.(*Link); ok {
		paramsMap = convertLink(*m)
	} else if m, ok := message.(Markdown); ok {
		paramsMap = convertMarkdown(m)
	} else if m, ok := message.(*Markdown); ok {
		paramsMap = convertMarkdown(*m)
	} else {
		return Result{ErrMsg: "not support message type"}
	}

	var buf []byte
	if buf, err = json.Marshal(paramsMap); err != nil {
		return Result{ErrMsg: "marshal message error:" + err.Error()}
	}

	return postMessage(DingAPIURL+ding.AccessToken, string(buf))
}

func convertMessage(m Message) map[string]interface{} {
	var paramsMap = make(map[string]interface{})
	paramsMap["msgtype"] = "text"
	paramsMap["text"] = map[string]string{"content": m.Content}
	paramsMap["at"] = map[string]interface{}{"atMobiles": m.AtPersion, "isAtAll": m.AtAll}
	return paramsMap
}

func convertLink(m Link) map[string]interface{} {
	var paramsMap = make(map[string]interface{})
	paramsMap["msgtype"] = "link"
	paramsMap["link"] = map[string]string{"text": m.Content, "title": m.Title, "picUrl": m.PictureURL, "messageUrl": m.MessageURL}
	return paramsMap
}

func convertMarkdown(m Markdown) map[string]interface{} {
	var paramsMap = make(map[string]interface{})
	paramsMap["msgtype"] = "markdown"
	paramsMap["markdown"] = map[string]string{"text": m.Content, "title": m.Title}
	return paramsMap
}

func postMessage(url string, message string) Result {
	var result Result

	resp, err := http.Post(url, "application/json", strings.NewReader(message))
	if err != nil {
		result.ErrMsg = "post data to api error:" + err.Error()
		return result
	}

	log.Println("message:", message)

	defer resp.Body.Close()
	var content []byte
	if content, err = ioutil.ReadAll(resp.Body); err != nil {
		result.ErrMsg = "read http response body error:" + err.Error()
		return result
	}

	log.Println("response result:", string(content))
	if err = json.Unmarshal(content, &result); err != nil {
		result.ErrMsg = "unmarshal http response body error:" + err.Error()
		return result
	}

	if result.ErrCode == 0 {
		result.Success = true
	}

	return result
}
