package bilibili

import (
	"context"
	"fmt"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/ssdomei232/goodBaby/internal/retry"
)

// SendDynamic 用给定账号发送一条 B 站动态
//
// 暂时没有处理 429 和 403 的区别对待
func SendDynamic(ctx context.Context, account *BiliAccount, config *BiliDynamicConfig) error {
	client := newClient(account.RawCookies)
	params := bilibili.CreateDynamicParam{
		DynamicId: 0,
		Type:      4,
		Rid:       0,
		Content:   config.Msg,
	}

	if err := retry.Do(ctx, func() error {
		_, err := client.CreateDynamic(params)
		return err
	}); err != nil {
		return fmt.Errorf("发送B站动态失败: %w", err)
	}

	return nil
}
