package rainyun

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/internal/meta"
)

// Rainyun规则验证器
type RainyunWorkorderRuleValidator struct{}

func (v *RainyunWorkorderRuleValidator) GetType() string {
	return RuleType
}

func (v *RainyunWorkorderRuleValidator) Validate(configJSON string) error {
	var config RainyunWorkOrderRule
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析雨云规则配置失败: %v", err)
	}

	if config.Title == "" {
		return fmt.Errorf("雨云规则配置中 title 不能为空")
	}

	if config.Msg == "" {
		return fmt.Errorf("雨云规则配置中 msg 不能为空")
	}

	return nil
}

func (v *RainyunWorkorderRuleValidator) Meta() meta.RuleMeta {
	return meta.RuleMeta{
		Type:        RuleType,
		Label:       "发送雨云工单",
		Description: "触发时通过雨云向指定的工单系统发送工单。",
		AccountType: AccountType,
		Fields: []meta.Field{
			{Key: "title", Label: "工单标题", Type: meta.FieldTextarea, Required: true},
			{Key: "msg", Label: "工单内容", Type: meta.FieldTextarea, Required: true},
		},
	}
}
