package bilibili

import (
	"context"
	"fmt"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

// SendBiliDynamicMsg 在 Bilibili 动态发送消息
//
// 暂时没有处理 429 和 403 的区别对待
func SendBiliDynamicMsg(ctx context.Context, rule *model.Rule) error {
	biliClient, err := getBiliClient(rule)
	if err != nil {
		return fmt.Errorf("获取B站客户端失败: %w", err)
	}

	biliDynamicConfig, err := getBiliDynamicConfig(rule)
	if err != nil {
		return fmt.Errorf("获取B站动态配置失败: %w", err)
	}

	dynamicParams := bilibili.CreateDynamicParam{
		DynamicId: 0,
		Type:      4,
		Rid:       0,
		Content:   biliDynamicConfig.Msg,
	}

	if err := retry.Do(ctx, func() error {
		_, err := biliClient.CreateDynamic(dynamicParams)
		return err
	}); err != nil {
		return fmt.Errorf("发送B站动态失败: %w", err)
	}

	return nil
}
