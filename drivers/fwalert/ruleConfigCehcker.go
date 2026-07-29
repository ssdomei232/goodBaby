package fwalert

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ssdomei232/goodBaby/internal/meta"
)

type FwalertRuleValidator struct{}

func (v *FwalertRuleValidator) GetType() string {
	return RuleTypeFwalert
}

func (v *FwalertRuleValidator) Validate(configJSON string) error {
	var config FwAlertRuleConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析饭碗警告规则配置失败: %v", err)
	}

	parsed, err := url.Parse(config.WebhookURL)
	if err != nil {
		return fmt.Errorf("饭碗警告规则配置中 webhook_url 格式不正确: %v", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("饭碗警告规则配置中 webhook_url 必须是 HTTP 或 HTTPS URL")
	}

	if config.Msg == "" {
		return fmt.Errorf("饭碗警告规则配置中 msg 不能为空")
	}
	if len([]rune(config.Msg)) > 1000 {
		return fmt.Errorf("饭碗警告内容过长")
	}

	return nil
}

func (v *FwalertRuleValidator) Meta() meta.RuleMeta {
	return meta.RuleMeta{
		Type:        RuleTypeFwalert,
		Label:       "发送饭碗警告",
		Description: "触发时发送饭碗警告消息。",
		Docs:        "docs/fwalert-config.md",
		AccountType: "",
		Fields: []meta.Field{
			{Key: "webhook_url", Label: "Webhook URL", Type: meta.FieldString, Required: true, Help: "饭碗警告的触发地址", Placeholder: "https://fwalert.com/xxx-xxx-xxx-xxx-xxx"},
			{Key: "msg", Label: "消息内容", Type: meta.FieldTextarea, Required: true, Help: "最大长度 1000 个字符"},
		},
	}
}

// ParseRuleConfig 解析饭碗警告规则配置
func ParseRuleConfig(configJSON string) (*FwAlertRuleConfig, error) {
	var config FwAlertRuleConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析饭碗警告规则配置失败: %v", err)
	}
	return &config, nil
}
