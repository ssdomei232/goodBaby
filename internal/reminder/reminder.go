package reminder

import (
	"fmt"
	"log"
	"time"

	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/internal/sender"
)

func Reminder() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
	}

	remainingTime := GlobalTimerManager.GetRemainingTime()
	if remainingTime <= 48*time.Hour {
		sendMsg := fmt.Sprintf("距离摇篮系统触发还有%d小时,请尽快发送Signal", int(remainingTime.Hours()))
		sender.SendDingTalkMsg("摇篮系统触发提醒", sendMsg, config.DingtalkBot.AccessToken, config.DingtalkBot.Secret)
	}
}
