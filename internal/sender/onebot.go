package sender

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/internal/base"
	napcat "github.com/ssdomei232/napcat-http-go-sdk"
)

func SendOneBot() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
		return
	}

	stopMsg := fmt.Sprintf("%s\n\n%s\n\n本消息由自动程序发送(摇篮系统)", base.GetBasicInfo(), config.OneBotConfig.Msg)

	for _, groupId := range config.OneBotConfig.SendGroups {
		sendOneBotMsg(groupId, stopMsg)
	}
}

// 发送一条消息到OneBot
func sendOneBotMsg(groupID int, msg string) {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
		return
	}

	client := napcat.NewClient(config.OneBotConfig.URL, config.OneBotConfig.Token)

	err = client.SendGroupMsg(strconv.Itoa(groupID), msg)
	if err != nil {
		log.Printf("发送消息到OneBot失败: %v", err)
	} else {
		log.Printf("成功发送消息到OneBot群 %d", groupID)
	}
}
