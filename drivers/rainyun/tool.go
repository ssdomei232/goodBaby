package rainyun

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// 从 Rule 中获取 RainyunConfig
func getRainyunConfigFromRule(rule *model.Rule) (*RainyunWorkOrderRule, error) {
	var rainyunConfig RainyunWorkOrderRule
	if err := json.Unmarshal([]byte(rule.ConfigJson), &rainyunConfig); err != nil {
		return nil, fmt.Errorf("解析Rainyun规则配置失败: %w", err)
	}
	return &rainyunConfig, nil
}

// 从 Rule 中获取 RainyunAccount
func getRainyunAccountFromRule(rule *model.Rule) (*RainyunAccount, error) {
	var rainyunAccount RainyunAccount
	if err := db.LoadAccountConfig(rule.AccountID, &rainyunAccount); err != nil {
		return nil, err
	}
	return &rainyunAccount, nil
}
