package dingtalk

import (
	"encoding/json"
	"fmt"
)

// ParseAccountConfig 解析钉钉机器人凭据配置
func ParseAccountConfig(config string) (*DingTalkConfig, error) {
	if config == "" {
		return nil, fmt.Errorf("解析钉钉机器人配置失败: 配置为空")
	}

	var account DingTalkConfig
	if err := json.Unmarshal([]byte(config), &account); err != nil {
		return nil, fmt.Errorf("解析钉钉机器人配置失败: %v", err)
	}
	return &account, nil
}

// ParseRuleConfig 解析钉钉规则配置
func ParseRuleConfig(configJSON string) (*DingTalkRuleConfig, error) {
	var config DingTalkRuleConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析钉钉规则配置失败: %v", err)
	}
	return &config, nil
}
