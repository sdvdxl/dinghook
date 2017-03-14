package dinghook

import (
	"fmt"
	"testing"
	"time"
)

func sTestSend(t *testing.T) {
	ding := Ding{AccessToken: ""}
	msg := Message{Content: "测试"}
	result := ding.Send(msg)
	fmt.Println(result)

	link := Link{Title: "link测试", Content: "测试", ContentURL: "https://www.baidu.com"}
	result = ding.Send(link)
	fmt.Println(result)

	markdown := Markdown{Title: "markdown测试", Content: "#### 杭州天气\n" +
		"> 9度，西北风1级，空气良89，相对温度73%\n\n" +
		"> ![screenshot](http://image.jpg)\n" +
		"> ###### 10点20分发布 [天气](http://www.thinkpage.cn/) \n"}
	result = ding.Send(markdown)
	fmt.Println(result)
}

func TestDingQueue(t *testing.T) {
	ding := &DingQueue{Title: "queue测试", Interval: 3, AccessToken: "91b35169899bc96e9648b2b8f4208ca56b6f84e14b137ba2b178e0cde9453817"}
	ding.Init()

	go ding.Start()

	ding.Push("#### 杭州天气\n" +
		"> 9度，西北风1级，空气良89，相对温度73%\n\n" +
		"> ![screenshot](http://image.jpg)\n" +
		"> ###### 10点20分发布 [天气](http://www.thinkpage.cn/) \n")
	time.Sleep(time.Second * 5)
	ding.Push("#### 杭州天气\n" +
		"> 9度，西北风1级，空气良89，相对温度73%\n\n" +
		"> ![screenshot](http://image.jpg)\n" +
		"> ###### 10点20分发布 [天气](http://www.thinkpage.cn/) \n")
	ding.Push("#### 杭州天气\n" +
		"> 9度，西北风1级，空气良89，相对温度73%\n\n" +
		"> ![screenshot](http://image.jpg)\n" +
		"> ###### 10点20分发布 [天气](http://www.thinkpage.cn/) \n")

	time.Sleep(time.Second * 10)

	ding.Push("#### 杭州天气\n" +
		"> 9度，西北风1级，空气良89，相对温度73%\n\n" +
		"> ![screenshot](http://image.jpg)\n" +
		"> ###### 10点20分发布 [天气](http://www.thinkpage.cn/) \n")
	ding.Push("#### 杭州天气\n" +
		"> 9度，西北风1级，空气良89，相对温度73%\n\n" +
		"> ![screenshot](http://image.jpg)\n" +
		"> ###### 10点20分发布 [天气](http://www.thinkpage.cn/) \n")

	time.Sleep(time.Second * 10)
}
