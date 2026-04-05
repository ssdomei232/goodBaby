package onebot

import (
	"context"
	"log"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/model"
)

func SendOneBotMsg(rule *model.Rule) {
	oneBotConfig := getOneBotConfigFromRule(rule)
	oneBotAccount, err := getOneBotAccountFromRule(rule)
	if err != nil {
		log.Printf("获取OneBot账户失败: %v", err)
		return
	}

	for _, groupID := range oneBotConfig.SendGroups {
		config, err := configs.GetConfig()
		if err != nil {
			log.Printf("获取配置失败: %v", err)
			return
		}
		var timeout time.Duration = time.Duration(config.TimeoutDurationHours) * time.Hour

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		operation := func() (string, error) {
			err := sendOneBotMsg(oneBotAccount, groupID, oneBotConfig.Msg)
			return "", err
		}

		_, err = backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
		if err != nil {
			log.Printf("发送OneBot消息失败: %v", err)
		}
	}
}
