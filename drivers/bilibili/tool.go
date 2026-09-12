package bilibili

import (
	"encoding/json"
	"fmt"

	"github.com/CuteReimu/bilibili/v2"
)

// ParseAccountConfig 解析 B 站账号配置
func ParseAccountConfig(config string) (*BiliAccount, error) {
	if config == "" {
		return nil, fmt.Errorf("解析B站账号配置失败: 规则没有关联账号")
	}

	var account BiliAccount
	if err := json.Unmarshal([]byte(config), &account); err != nil {
		return nil, fmt.Errorf("解析B站账号配置失败: %v", err)
	}
	return &account, nil
}

// ParseDynamicConfig 解析 B 站动态规则配置
func ParseDynamicConfig(configJSON string) (*BiliDynamicConfig, error) {
	var config BiliDynamicConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析B站动态规则配置失败: %v", err)
	}
	return &config, nil
}

// ParsePrivateMessageConfig 解析 B 站私信规则配置
func ParsePrivateMessageConfig(configJSON string) (*BiliPrivateMessageConfig, error) {
	var config BiliPrivateMessageConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析B站私信规则配置失败: %v", err)
	}
	return &config, nil
}

// newClient 按账号配置创建 B 站客户端
func newClient(rawCookies string) *bilibili.Client {
	client := bilibili.New()
	client.SetRawCookies(rawCookies)
	return client
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
