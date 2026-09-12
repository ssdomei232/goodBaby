package onebot

import (
	"context"
	"fmt"
	"strings"

	"github.com/ssdomei232/goodBaby/internal/retry"
)

// SendMessage 向规则中配置的群和好友发送消息
//
// 每个目标单独重试、互不阻塞，最后汇总失败信息交给上层记录。
func SendMessage(ctx context.Context, account *OneBotAccount, config *OneBotConfig) error {
	var fails []string
	total := 0

	for _, groupID := range config.SendGroups {
		total++
		if err := retry.Do(ctx, func() error {
			return sendGroupMsg(account, groupID, config.Msg)
		}); err != nil {
			fails = append(fails, fmt.Sprintf("群 %d: %v", groupID, err))
		}
	}

	for _, userID := range config.SendUsers {
		total++
		if err := retry.Do(ctx, func() error {
			return sendPrivateMsg(account, userID, config.Msg)
		}); err != nil {
			fails = append(fails, fmt.Sprintf("好友 %d: %v", userID, err))
		}
	}

	if len(fails) > 0 {
		return fmt.Errorf("%d/%d 条消息发送失败: %s", len(fails), total, strings.Join(fails, "; "))
	}
	return nil
}
