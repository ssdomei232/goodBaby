// Package dingtalk 提供钉钉自定义机器人的提醒与规则执行能力
package dingtalk

import (
	"fmt"

	"github.com/CatchZeng/dingtalk/pkg/dingtalk"
	"github.com/ssdomei232/goodBaby/model"
)

// SendDingTalkMsg 使用用户配置的钉钉机器人发送提醒
//
// 用户没有配置钉钉机器人时返回 ErrNoUserConfig
func SendDingTalkMsg(user *model.User, title string, msg string) error {
	config := getDingTalkConfigFromUser(user)
	if config == nil {
		return ErrNoUserConfig
	}
	return send(config, title, msg)
}

// ErrNoUserConfig 表示用户没有配置钉钉提醒
var ErrNoUserConfig = fmt.Errorf("用户未配置钉钉机器人")

func send(config *DingTalkConfig, title, msg string) error {
	client := dingtalk.NewClient(config.AccessToken, config.Secret)
	message := dingtalk.NewMarkdownMessage().SetMarkdown(title, msg)

	if _, _, err := client.Send(message); err != nil {
		return fmt.Errorf("发送钉钉消息失败: %w", err)
	}
	return nil
}
