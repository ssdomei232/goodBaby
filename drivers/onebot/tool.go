package onebot

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
	napcat "github.com/ssdomei232/napcat-http-go-sdk"
)

func newClient(account *OneBotAccount) *napcat.Client {
	return napcat.NewClient(account.Token, account.URL)
}

// 从 Rule 中获取 OneBotConfig
func getOneBotConfigFromRule(rule *model.Rule) (*OneBotConfig, error) {
	var oneBotConfig OneBotConfig
	if err := json.Unmarshal([]byte(rule.ConfigJson), &oneBotConfig); err != nil {
		return nil, fmt.Errorf("解析OneBot规则配置失败: %w", err)
	}
	return &oneBotConfig, nil
}

// 从 Rule 中获取 OneBotAccount
func getOneBotAccountFromRule(rule *model.Rule) (*OneBotAccount, error) {
	var oneBotAccount OneBotAccount
	if err := db.LoadAccountConfig(rule.AccountID, &oneBotAccount); err != nil {
		return nil, err
	}
	return &oneBotAccount, nil
}

// 发送一条群消息
//
// SDK 只返回传输层错误，这里自己解析 OneBot 的响应信封，
// 否则 token 错误、群号不存在这类失败会被当成发送成功。
func sendGroupMsg(account *OneBotAccount, groupID int64, msg string) error {
	return call(account, "/send_group_msg", map[string]any{
		"group_id": groupID,
		"message":  msg,
	})
}

// 发送一条私聊消息
func sendPrivateMsg(account *OneBotAccount, userID int64, msg string) error {
	return call(account, "/send_private_msg", map[string]any{
		"user_id": userID,
		"message": msg,
	})
}

func call(account *OneBotAccount, endpoint string, payload map[string]any) error {
	var resp apiResponse
	if err := newClient(account).DoRequest("POST", endpoint, payload, &resp); err != nil {
		return fmt.Errorf("请求 OneBot 失败: %w", err)
	}
	return resp.err()
}
