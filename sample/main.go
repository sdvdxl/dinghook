package main

import (
	"fmt"
	"os"

	"github.com/sdvdxl/dinghook"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("program token [signSecret]")
		return
	}

	ding := dinghook.NewDing(os.Args[1])
	ding.SignToken = os.Args[2]
	// fmt.Println(ding.Send(dinghook.Message{Content: "内容"}))

	// fmt.Println(ding.Send(dinghook.Markdown{Title: "标题", Content: "# 标题 \n## 二级标题"}))

	// fmt.Println(ding.Send(dinghook.Link{Title: "标题", Content: "内容简介", ContentURL: "https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq", PictureURL: "https://dingtalkdoc.oss-cn-beijing.aliyuncs.com/images/0.0.204/1571983069016-873ac5a2-fc5c-4281-bf85-48d02e05b9b6.png"}))

	// fmt.Println(ding.Send(dinghook.OverallActionCard{Title: "标题", Content: "![](https://img.alicdn.com/tfs/TB1nhWCiBfH8KJjy1XbXXbLdXXa-547-379.png)\n# 内容简介", ButtonURL: "https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq", ButtonTitle: "查看"}))

	// fmt.Println(ding.Send(dinghook.IndependentActionCard{Title: "标题", Content: "![](https://img.alicdn.com/tfs/TB1nhWCiBfH8KJjy1XbXXbLdXXa-547-379.png)\n# 内容简介",
	// 	Btns: []dinghook.IndependentActionCardButton{
	// 		{ButtonTitle: "测试1", ButtonURL: "https://www.baidu.com"},
	// 		{ButtonTitle: "测试2", ButtonURL: "https://www.weibo.com"},
	// 	}}))

	fmt.Println(ding.Send(dinghook.FeedCard{Links: []dinghook.FeedCardLink{
		{Title: "测试", ContentURL: "https://www.baidu.com", PictureURL: "https://dingtalkdoc.oss-cn-beijing.aliyuncs.com/images/0.0.204/1571983069016-873ac5a2-fc5c-4281-bf85-48d02e05b9b6.png"},
		{Title: "测试", ContentURL: "https://www.baidu.com", PictureURL: "https://dingtalkdoc.oss-cn-beijing.aliyuncs.com/images/0.0.204/1571983069016-873ac5a2-fc5c-4281-bf85-48d02e05b9b6.png"},
	}}))
}
