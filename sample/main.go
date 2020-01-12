package main

import (
	"fmt"

	"github.com/sdvdxl/dinghook"
)

func main() {
	ding := dinghook.NewDing("token")
	ding.SignToken = "sign secret"
	// text 类型消息
	fmt.Println(ding.Send(dinghook.Message{Content: "内容"}))
}
