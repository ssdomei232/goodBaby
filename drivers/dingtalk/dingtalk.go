// dingtalk 仅用于发送提醒
package dingtalk

import (
	"log"

	"github.com/CatchZeng/dingtalk/pkg/dingtalk"
	"github.com/ssdomei232/goodBaby/model"
)

func SendDingTalkMsg(user *model.User, title string, msg string) {
	dingtalkAccount := getDingTalkConfigFromUser(user)
	if dingtalkAccount == nil {
		return
	}

	client := dingtalk.NewClient(dingtalkAccount.AccessToken, dingtalkAccount.Secret)
	sendMsg := dingtalk.NewMarkdownMessage().SetMarkdown(title, msg)
	_, _, err := client.Send(sendMsg)
	if err != nil {
		log.Printf("发送钉钉消息失败: %v", err)
	}
}
