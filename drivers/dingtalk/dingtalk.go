// Package dingtalk 提供钉钉自定义机器人的提醒与规则执行能力
package dingtalk

import (
	"context"
	"fmt"

	"github.com/CatchZeng/dingtalk/pkg/dingtalk"
	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

// ErrNoUserConfig 表示用户没有配置钉钉提醒
var ErrNoUserConfig = fmt.Errorf("用户未配置钉钉机器人")

// SendDingTalkMsg 使用用户配置的钉钉机器人发送提醒
//
// 用户没有配置钉钉机器人时返回 ErrNoUserConfig
func SendDingTalkMsg(user *model.User, title string, msg string) error {
	if user.DingTalkConfig == nil {
		return ErrNoUserConfig
	}

	config, err := ParseAccountConfig(*user.DingTalkConfig)
	if err != nil {
		return err
	}

	if err := send(config, title, msg); err != nil {
		return fmt.Errorf("发送钉钉消息失败: %w", err)
	}
	return nil
}

// SendMessage 用规则里配置的机器人发送消息
func SendMessage(ctx context.Context, config *DingTalkRuleConfig) error {
	if err := retry.Do(ctx, func() error {
		return send(&config.DingTalkConfig, config.Title, config.Msg)
	}); err != nil {
		return fmt.Errorf("发送钉钉消息失败: %w", err)
	}
	return nil
}

// send 调用一次机器人接口
func send(config *DingTalkConfig, title, msg string) error {
	client := dingtalk.NewClient(config.AccessToken, config.Secret)
	message := dingtalk.NewMarkdownMessage().SetMarkdown(title, msg)

	_, _, err := client.Send(message)
	return err
}
