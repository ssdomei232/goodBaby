package onebot

import (
	"context"
	"fmt"
	"strings"

	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

// SendOneBotMsg 向规则中配置的群和好友发送消息
func SendOneBotMsg(ctx context.Context, rule *model.Rule) error {
	oneBotConfig, err := getOneBotConfigFromRule(rule)
	if err != nil {
		return err
	}

	oneBotAccount, err := getOneBotAccountFromRule(rule)
	if err != nil {
		return fmt.Errorf("获取OneBot账户失败: %w", err)
	}

	var fails []string
	total := 0

	for _, groupID := range oneBotConfig.SendGroups {
		total++
		err := retry.Do(ctx, func() error {
			return sendGroupMsg(oneBotAccount, groupID, oneBotConfig.Msg)
		})
		if err != nil {
			fails = append(fails, fmt.Sprintf("群 %d: %v", groupID, err))
		}
	}

	for _, userID := range oneBotConfig.SendUsers {
		total++
		err := retry.Do(ctx, func() error {
			return sendPrivateMsg(oneBotAccount, userID, oneBotConfig.Msg)
		})
		if err != nil {
			fails = append(fails, fmt.Sprintf("好友 %d: %v", userID, err))
		}
	}

	if len(fails) > 0 {
		return fmt.Errorf("%d/%d 条消息发送失败: %s", len(fails), total, strings.Join(fails, "; "))
	}
	return nil
}
