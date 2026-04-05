package bilibili

import (
	"context"
	"log"
	"time"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/cenkalti/backoff/v5"
	"github.com/ssdomei232/goodBaby/configs"
)

// 在Bilibili动态发送消息
//
// 暂时没有处理 429 和 403 的区别对待
func SendBiliDynamicMsg(biliClient *bilibili.Client, msg string) {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置失败: %v", err)
		return
	}
	var timeout time.Duration = time.Duration(config.TimeoutDurationHours) * time.Hour

	var dynamicParams bilibili.CreateDynamicParam
	dynamicParams = bilibili.CreateDynamicParam{
		DynamicId: 0,
		Type:      4,
		Rid:       0,
		Content:   msg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	operation := func() (string, error) {
		_, err := biliClient.CreateDynamic(dynamicParams)
		return "", err
	}

	_, err = backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		log.Printf("发送B站动态失败: %v", err)
	}
}
