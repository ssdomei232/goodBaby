package email

import (
	"encoding/json"
	"fmt"
	"net/mail"

	"github.com/ssdomei232/goodBaby/internal/meta"
)

// EmailRuleValidator 邮件规则验证器
type EmailRuleValidator struct{}

func (v *EmailRuleValidator) GetType() string {
	return RuleType
}

func (v *EmailRuleValidator) Validate(configJSON string) error {
	var config EmailRule
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析邮件规则配置失败: %v", err)
	}

	if config.Title == "" {
		return fmt.Errorf("邮件规则配置中 title 不能为空")
	}

	if config.Msg == "" {
		return fmt.Errorf("邮件规则配置中 msg 不能为空")
	}

	if len(config.Destinations) == 0 {
		return fmt.Errorf("邮件规则配置中 destinations 不能为空")
	}

	// 验证邮箱地址格式
	for _, dest := range config.Destinations {
		if _, err := mail.ParseAddress(dest); err != nil {
			return fmt.Errorf("无效的邮箱地址: %s", dest)
		}
	}

	return nil
}

func (v *EmailRuleValidator) Meta() meta.RuleMeta {
	return meta.RuleMeta{
		Type:        RuleType,
		Label:       "发送邮件",
		Description: "触发时通过关联的 SMTP 账号向指定地址发送邮件。",
		Docs:        "docs/email-config.md",
		AccountType: AccountType,
		Fields: []meta.Field{
			{Key: "title", Label: "邮件标题", Type: meta.FieldString, Required: true},
			{Key: "msg", Label: "邮件正文", Type: meta.FieldTextarea, Required: true},
			{Key: "destinations", Label: "收件人", Type: meta.FieldStringList, Required: true, Placeholder: "someone@example.com"},
		},
	}
}
