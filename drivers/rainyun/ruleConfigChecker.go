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
		Description: "触发时发送雨云工单",
		AccountType: AccountType,
		Fields: []meta.Field{
			{Key: "title", Label: "工单标题", Type: meta.FieldTextarea, Required: true},
			{Key: "msg", Label: "工单内容", Type: meta.FieldTextarea, Required: true},
		},
	}
}

type RainyunRunAwayRuleValidator struct{}

func (v *RainyunRunAwayRuleValidator) GetType() string {
	return "rainyun-runaway"
}

func (v *RainyunRunAwayRuleValidator) Validate(configJSON string) error {
	var config RainyunRunAwayRule
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析雨云规则配置失败: %v", err)
	}

	if config.IK != "我已知晓" {
		return fmt.Errorf("请输入正确内容")
	}

	return nil
}

func (v *RainyunRunAwayRuleValidator) Meta() meta.RuleMeta {
	return meta.RuleMeta{
		Type:        "rainyun-runaway",
		Label:       "重置雨云账号中所有云服务器(一键跑路)",
		Description: "触发时会重装雨云账号下所有云服务器来实现跑路",
		AccountType: AccountType,
		Fields: []meta.Field{
			{Key: "ik", Label: "该操作引发的损失及法律责任归属账号所有人，本软件一律不承担，请输入我已知晓", Type: meta.FieldTextarea, Required: true},
		},
	}
}
