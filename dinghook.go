// Copyright 2020 sdvdxl
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dinghook

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"container/list"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"sync"
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

// Group 钉钉组
type Group struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

// Ding 钉钉消息发送实体
type Ding struct {
	AccessToken string // token
	SignToken   string // 加签token，可选，如果填写则使用加签方式发送，否则需要自己使用关键词或者ip方式
}

func calcSign(timestamp int64, signToken string) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, signToken)
	h := hmac.New(sha256.New, []byte(signToken))
	h.Write([]byte(stringToSign))
	sumValue := h.Sum(nil)
	return url.QueryEscape(base64.StdEncoding.EncodeToString(sumValue))
}

// NewDing new 一个没有队列的ding
func NewDing(token string) Ding {
	return Ding{AccessToken: token}
}

// DingQueue 用queue 方式发送消息
// 会发送 markdown 类型消息
type DingQueue struct {
	AccessToken string
	ding        Ding
	Interval    uint       // 发送间隔s，最小为1
	Limit       uint       // 每次发送消息限制，0 无限制，到达时间则发送队列所有消息，大于1则到时间发送最大Limit数量的消息
	Title       string     // 摘要
	messages    *list.List // 消息队列
	lock        sync.Mutex
}

// NewQueue 创建一个队列
func NewQueue(token, title string, interval, limit uint) *DingQueue {
	dingQueue := &DingQueue{
		AccessToken: token,
		Title:       title,
		Interval:    interval,
		Limit:       limit}
	dingQueue.Init()
	return dingQueue
}

// Init 初始化 DingQueue
func (ding *DingQueue) Init() {
	ding.ding.AccessToken = ding.AccessToken
	ding.messages = list.New()
	if ding.Interval == 0 {
		ding.Interval = 1
	}
}

// Push push 消息到队列
func (ding *DingQueue) Push(message string) {
	defer ding.lock.Unlock()
	ding.lock.Lock()
	ding.messages.PushBack(SimpleMessage{Title: ding.Title, Content: message})
}

// PushWithTitle push 消息到队列
func (ding *DingQueue) PushWithTitle(title, message string) {
	defer ding.lock.Unlock()
	ding.lock.Lock()
	if title == "" {
		title = ding.Title
	}

	ding.messages.PushBack(SimpleMessage{Title: title, Content: message})
}

// PushMessage push 消息到队列
func (ding *DingQueue) PushMessage(m SimpleMessage) {
	defer ding.lock.Unlock()
	ding.lock.Lock()
	ding.messages.PushBack(m)
}

// Start 开始工作
func (ding *DingQueue) Start() {
	sendQueueMessage(ding)
	timer := time.NewTicker(time.Second * time.Duration(ding.Interval))
	for {
		select {
		case <-timer.C:
			sendQueueMessage(ding)
		}
	}
}

func sendQueueMessage(ding *DingQueue) {
	defer ding.lock.Unlock()
	ding.lock.Lock()
	title := ding.Title
	msg := ""
	if ding.Limit == 0 {
		for {
			m := ding.messages.Front()
			if m == nil {
				break
			}
			ding.messages.Remove(m)
			switch m.Value.(type) {
			case SimpleMessage:
				v := m.Value.(SimpleMessage)
				msg += v.Content + "\n\n"

			case string:
				msg += m.Value.(string) + "\n\n"
			}

		}
	} else {
	label:
		for i := uint(0); i < ding.Limit; i++ {
			for {
				m := ding.messages.Front()

				if m == nil {
					break label
				}
				ding.messages.Remove(m)
				switch m.Value.(type) {
				case SimpleMessage:
					v := m.Value.(SimpleMessage)
					msg += v.Content + "\n\n"
				case string:
					msg += m.Value.(string) + "\n\n"
				}
			}
		}
	}

	if msg != "" {
		go func() {
			r := ding.ding.Send(Markdown{Title: title, Content: msg})
			if !r.Success {
				log.Println("err:" + r.ErrMsg)
				NewDing(ding.ding.AccessToken).Send("消息太长，请通过其他途径查看，比如邮件")
			}
		}()
	}
}

// SendMessage 发送普通文本消息
func (ding Ding) SendMessage(message Message) Result {
	return ding.Send(message)
}

// SendLink 发送link类型消息
func (ding Ding) SendLink(message Link) Result {
	return ding.Send(message)
}

// SendMarkdown 发送markdown格式消息
func (ding Ding) SendMarkdown(message Markdown) Result {
	return ding.Send(message)
}

// Send 发送消息
func (ding Ding) Send(message interface{}) Result {
	if ding.AccessToken == "" {
		return Result{ErrMsg: "access token is required"}
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
	} else if m, ok := message.(OverallActionCard); ok {
		paramsMap = convertOverallActionCard(m)
	} else if m, ok := message.(*OverallActionCard); ok {
		paramsMap = convertOverallActionCard(*m)
	} else if m, ok := message.(IndependentActionCard); ok {
		paramsMap = convertIndependentActionCard(m)
	} else if m, ok := message.(*IndependentActionCard); ok {
		paramsMap = convertIndependentActionCard(*m)
	} else if m, ok := message.(FeedCard); ok {
		paramsMap = convertFeedCard(m)
	} else if m, ok := message.(*FeedCard); ok {
		paramsMap = convertFeedCard(*m)
	} else {
		return Result{ErrMsg: "not support message type"}
	}

	var buf []byte
	if buf, err = json.Marshal(paramsMap); err != nil {
		return Result{ErrMsg: "marshal message error:" + err.Error()}
	}

	dingUrl := DingAPIURL + ding.AccessToken
	if ding.SignToken != "" {
		timestamp := time.Now().UnixNano() / 1000 / 1000
		dingUrl += fmt.Sprintf("&timestamp=%d&sign=%s", timestamp, calcSign(timestamp, ding.SignToken))
	}
	return postMessage(dingUrl, string(buf))
}

func convertMessage(m Message) map[string]interface{} {
	var paramsMap = make(map[string]interface{})
	paramsMap["msgtype"] = "text"
	paramsMap["text"] = map[string]string{"content": m.Content}
	paramsMap["at"] = map[string]interface{}{"atMobiles": m.AtData.AtMobiles,"atUserIds":m.AtData, "isAtAll": m.AtData}
	return paramsMap
}

func convertLink(m Link) map[string]interface{} {
	var paramsMap = make(map[string]interface{})
	paramsMap["msgtype"] = "link"
	paramsMap["link"] = map[string]string{"text": m.Content, "title": m.Title, "picUrl": m.PictureURL, "messageUrl": m.ContentURL}
	return paramsMap
}

func convertMarkdown(m Markdown) map[string]interface{} {
	var paramsMap = make(map[string]interface{})
	paramsMap["msgtype"] = "markdown"
	paramsMap["markdown"] = map[string]string{"text": m.Content, "title": m.Title}
	return paramsMap
}

func convertOverallActionCard(m OverallActionCard) map[string]interface{} {
	var paramsMap = make(map[string]interface{})
	paramsMap["msgtype"] = "actionCard"

	btnOrientation := "0"
	if m.ButtonHorizontal {
		btnOrientation = "1"
	}

	paramsMap["actionCard"] = map[string]string{"text": m.Content, "title": m.Title,
		"singleTitle": m.ButtonTitle, "singleURL": m.ButtonURL,
		"btnOrientation": btnOrientation}
	return paramsMap
}

func convertIndependentActionCard(m IndependentActionCard) map[string]interface{} {
	var paramsMap = make(map[string]interface{})
	paramsMap["msgtype"] = "actionCard"

	btnOrientation := "0"
	if m.ButtonHorizontal {
		btnOrientation = "1"
	}

	btns := make([]map[string]interface{}, 0, len(m.Btns))

	for _, v := range m.Btns {
		btns = append(btns, map[string]interface{}{"title": v.ButtonTitle, "actionURL": v.ButtonURL})
	}

	paramsMap["actionCard"] = map[string]interface{}{"text": m.Content, "title": m.Title,
		"btns":           btns,
		"btnOrientation": btnOrientation}
	return paramsMap
}

func convertFeedCard(m FeedCard) map[string]interface{} {
	var paramsMap = make(map[string]interface{})
	paramsMap["msgtype"] = "feedCard"

	links := make([]map[string]interface{}, 0, len(m.Links))

	for _, v := range m.Links {
		links = append(links, map[string]interface{}{"title": v.Title, "messageURL": v.ContentURL, "picURL": v.PictureURL})
	}

	paramsMap["feedCard"] = map[string]interface{}{"links": links}
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
