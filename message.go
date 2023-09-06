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

//Package dinghook 详情参见 https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
package dinghook

const (
	// MsgTypeText text 类型
	MsgTypeText = "text"

	// MsgTypeLink link 类型
	MsgTypeLink = "link"

	// MsgTypeMarkdown markdown 类型
	MsgTypeMarkdown = "markdown"

	// MsgTypeActionCard actionCard 类型
	MsgTypeActionCard = "actionCard"
)

// At 结构
type At struct {
	AtMobiles []string
	AtUserIds []string
	AtAll     bool
}

// Message 普通消息
type Message struct {
	Content   string `validate:"required"`
	AtData At
}

// Link 链接消息
type Link struct {
	Content    string `json:"text" validate:"required"`       // 要发送的消息， 必填
	Title      string `json:"title" validate:"required"`      // 标题， 必填
	ContentURL string `json:"messageUrl" validate:"required"` // 点击消息跳转的URL 必填
	PictureURL string `json:"picUrl"`                         // 图片 url
}

// Markdown markdown 类型
type Markdown struct {
	Content string `json:"text" validate:"required"`  // 要发送的消息， 必填
	Title   string `json:"title" validate:"required"` // 标题， 必填
}

// OverallActionCard 整体跳转 类型
type OverallActionCard struct {
	Content     string `json:"text" validate:"required"`        // 要发送的消息， 必填
	Title       string `json:"title" validate:"required"`       // 标题， 必填
	ButtonHorizontal bool                          `json:"btnOrientation"`            //false-按钮竖直排列，true-按钮横向排列
	ButtonTitle string `json:"singleTitle" validate:"required"` // 单个按钮的文案
	ButtonURL   string `json:"singleURL" validate:"required"`   // 点击按钮跳转url
}

// IndependentActionCard 独立跳转类型
type IndependentActionCard struct {
	Content          string                        `json:"text" validate:"required"`  // 要发送的消息， 必填
	Title            string                        `json:"title" validate:"required"` // 标题， 必填
	ButtonHorizontal bool                          `json:"btnOrientation"`            //false-按钮竖直排列，true-按钮横向排列
	Btns             []IndependentActionCardButton `json:"btns" validate:"required"`
}

// IndependentActionCardButton 按钮信息
type IndependentActionCardButton struct {
	ButtonTitle string `json:"title" validate:"required"`     // 单个按钮的文案
	ButtonURL   string `json:"actionURL" validate:"required"` // 点击按钮跳转url
}

// FeedCard FeedCard类型
type FeedCard struct {
	Links []FeedCardLink
}

// FeedCardLink feedCard 链接
type FeedCardLink struct {
	Title      string `json:"title"  validate:"required"`
	ContentURL string `json:"messageUrl" validate:"required"` // 点击消息跳转的URL 必填
	PictureURL string `json:"picUrl"`                         // 图片 url
}

// SimpleMessage push message
type SimpleMessage struct {
	Content string
	Title   string
}
