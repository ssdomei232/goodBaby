package bilibili

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

// 私信消息类型：1 为文字
const privateMsgTypeText = 1

// SendBiliPrivateMessage 向配置中的每个 UID 发送一条 B 站私信
func SendBiliPrivateMessage(ctx context.Context, rule *model.Rule) error {
	client, err := getBiliClient(rule)
	if err != nil {
		return fmt.Errorf("获取B站客户端失败: %w", err)
	}

	config, err := getPrivateMessageConfig(rule)
	if err != nil {
		return err
	}

	// 发私信需要带上自己的 UID，从账号信息里取
	account, err := client.GetAccountInformation()
	if err != nil {
		return fmt.Errorf("获取B站账号信息失败(cookie 可能已失效): %w", err)
	}
	if account == nil || account.Mid == 0 {
		return fmt.Errorf("B站 cookie 无效或已过期")
	}

	var fails []string
	for _, uid := range config.ReceiverUids {
		err := retry.Do(ctx, func() error {
			_, err := client.SendPrivateMessage(bilibili.SendPrivateMessageParam{
				SenderUid:    account.Mid,
				ReceiverId:   int(uid),
				ReceiverType: 1,
				MsgType:      privateMsgTypeText,
				Timestamp:    int(time.Now().Unix()),
				Content:      buildTextContent(config.Msg),
			})
			return err
		})
		if err != nil {
			fails = append(fails, fmt.Sprintf("UID %d: %v", uid, err))
		}
	}

	if len(fails) > 0 {
		return fmt.Errorf("%d/%d 条私信发送失败: %s",
			len(fails), len(config.ReceiverUids), strings.Join(fails, "; "))
	}
	return nil
}

// buildTextContent 文字私信的 content 是一个 JSON 字符串 {"content":"..."}
func buildTextContent(msg string) string {
	payload, err := json.Marshal(map[string]string{"content": msg})
	if err != nil {
		// msg 是普通字符串，序列化不会失败；兜底也返回合法 JSON
		return `{"content":""}`
	}
	return string(payload)
}

func getPrivateMessageConfig(rule *model.Rule) (*BiliPrivateMessageConfig, error) {
	var config BiliPrivateMessageConfig
	if err := json.Unmarshal([]byte(rule.ConfigJson), &config); err != nil {
		return nil, fmt.Errorf("解析B站私信规则配置失败: %w", err)
	}
	return &config, nil
}
