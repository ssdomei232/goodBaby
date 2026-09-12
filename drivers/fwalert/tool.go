package fwalert

import (
	"encoding/json"
	"fmt"
)

// ParseRuleConfig 解析饭碗警告规则配置
func ParseRuleConfig(configJSON string) (*FwAlertRuleConfig, error) {
	var config FwAlertRuleConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析饭碗警告规则配置失败: %v", err)
	}
	return &config, nil
}
