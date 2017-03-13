package dinghook

import (
	"fmt"
	"testing"
)

func TestSend(t *testing.T) {
	ding := Ding{AccessToken: ""}
	msg := Message{Content: "测试"}
	result := ding.Send(msg)
	fmt.Println(result)
}
