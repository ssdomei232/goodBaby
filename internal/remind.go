package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/CatchZeng/dingtalk/pkg/dingtalk"
	"github.com/ssdomei232/goodBaby/configs"
)

func Reminder() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
	}

	remainingTime := GlobalTimerManager.GetRemainingTime()
	if remainingTime <= 48*time.Hour {
		sendMsg := fmt.Sprintf("距离摇篮系统触发还有%d小时,请尽快发送Signal", int(remainingTime.Hours()))
		sendDingTalkMsg("摇篮系统触发提醒", sendMsg, config.DingtalkBot.AccessToken, config.DingtalkBot.Secret)
	}
}

func sendDingTalkMsg(title string, msg string, accessToken string, secret string) {
	client := dingtalk.NewClient(accessToken, secret)
	sendMsg := dingtalk.NewMarkdownMessage().SetMarkdown(title, msg)
	_, _, err := client.Send(sendMsg)
	if err != nil {
		log.Printf("发送钉钉消息失败: %v", err)
	}
}
