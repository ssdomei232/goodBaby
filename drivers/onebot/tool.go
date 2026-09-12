package onebot

import (
	"encoding/json"
	"fmt"

	napcat "github.com/ssdomei232/napcat-http-go-sdk"
)

// ParseAccountConfig 解析 OneBot 账号配置
func ParseAccountConfig(config string) (*OneBotAccount, error) {
	if config == "" {
		return nil, fmt.Errorf("解析OneBot账号配置失败: 规则没有关联账号")
	}

	var account OneBotAccount
	if err := json.Unmarshal([]byte(config), &account); err != nil {
		return nil, fmt.Errorf("解析OneBot账号配置失败: %v", err)
	}
	return &account, nil
}

// ParseRuleConfig 解析 OneBot 规则配置
func ParseRuleConfig(configJSON string) (*OneBotConfig, error) {
	var config OneBotConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析OneBot规则配置失败: %v", err)
	}
	return &config, nil
}

// newClient 创建 OneBot HTTP 客户端
func newClient(account *OneBotAccount) *napcat.Client {
	return napcat.NewClient(account.Token, account.URL)
}

// sendGroupMsg 发送一条群消息
func sendGroupMsg(account *OneBotAccount, groupID int64, msg string) error {
	return call(account, "/send_group_msg", map[string]any{
		"group_id": groupID,
		"message":  msg,
	})
}

// sendPrivateMsg 发送一条私聊消息
func sendPrivateMsg(account *OneBotAccount, userID int64, msg string) error {
	return call(account, "/send_private_msg", map[string]any{
		"user_id": userID,
		"message": msg,
	})
}

// call 调用一次 OneBot HTTP API。
//
// SDK 只返回传输层错误，这里自己解析 OneBot 的响应信封，
// 否则 token 错误、群号不存在这类失败会被当成发送成功。
func call(account *OneBotAccount, endpoint string, payload map[string]any) error {
	var resp apiResponse
	if err := newClient(account).DoRequest("POST", endpoint, payload, &resp); err != nil {
		return fmt.Errorf("请求 OneBot 失败: %w", err)
	}
	return resp.err()
}
