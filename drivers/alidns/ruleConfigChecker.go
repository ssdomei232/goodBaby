package alidns

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/internal/meta"
)

// 阿里云删除记录规则验证器
type AliDNSDeleteRecordRuleValidator struct{}

func (v *AliDNSDeleteRecordRuleValidator) GetType() string {
	return RuleTypeDeleteRecord
}

func (v *AliDNSDeleteRecordRuleValidator) Validate(configJSON string) error {
	var config DeleteRecordConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析阿里云删除记录规则配置失败: %v", err)
	}

	if config.RecordID == "" {
		return fmt.Errorf("阿里云删除记录规则配置中 record_id 不能为空")
	}

	return nil
}

func (v *AliDNSDeleteRecordRuleValidator) Meta() meta.RuleMeta {
	return meta.RuleMeta{
		Type:        RuleTypeDeleteRecord,
		Label:       "删除阿里云DNS解析记录",
		Description: "触发时删除阿里云DNS解析记录",
		AccountType: AccountType,
		Fields: []meta.Field{
			{Key: "record_id", Label: "记录ID", Type: meta.FieldString, Required: true},
		},
	}
}
