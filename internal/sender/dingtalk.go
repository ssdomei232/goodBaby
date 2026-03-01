package sender

import (
	"log"

	"github.com/CatchZeng/dingtalk/pkg/dingtalk"
)

func SendDingTalkMsg(title string, msg string, accessToken string, secret string) {
	client := dingtalk.NewClient(accessToken, secret)
	sendMsg := dingtalk.NewMarkdownMessage().SetMarkdown(title, msg)
	_, _, err := client.Send(sendMsg)
	if err != nil {
		log.Printf("发送钉钉消息失败: %v", err)
	}
}
