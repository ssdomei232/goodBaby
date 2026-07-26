package dingtalk

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/internal/meta"
)

// DingTalkRuleValidator 钉钉规则验证器
type DingTalkRuleValidator struct{}

func (v *DingTalkRuleValidator) GetType() string {
	return RuleType
}

func (v *DingTalkRuleValidator) Validate(configJSON string) error {
	config, err := ParseRuleConfig(configJSON)
	if err != nil {
		return err
	}

	if config.AccessToken == "" {
		return fmt.Errorf("钉钉规则配置中 access_token 不能为空")
	}
	if config.Title == "" {
		return fmt.Errorf("钉钉规则配置中 title 不能为空")
	}
	if config.Msg == "" {
		return fmt.Errorf("钉钉规则配置中 msg 不能为空")
	}

	return nil
}

func (v *DingTalkRuleValidator) Meta() meta.RuleMeta {
	return meta.RuleMeta{
		Type:        RuleType,
		Label:       "钉钉机器人消息",
		Description: "触发时通过钉钉自定义机器人发送一条 Markdown 消息，不需要关联账号。",
		AccountType: "", // 凭据直接写在规则里，无需账号
		Fields: []meta.Field{
			{Key: "access_token", Label: "Access Token", Type: meta.FieldPassword, Required: true, Secret: true, Help: "机器人 Webhook 中 access_token 参数的值"},
			{Key: "secret", Label: "加签 Secret", Type: meta.FieldPassword, Secret: true, Help: "机器人安全设置选择“加签”时填写"},
			{Key: "title", Label: "消息标题", Type: meta.FieldString, Required: true},
			{Key: "msg", Label: "消息内容", Type: meta.FieldTextarea, Required: true, Help: "支持 Markdown"},
		},
	}
}

// ParseRuleConfig 解析钉钉规则配置
func ParseRuleConfig(configJSON string) (*DingTalkRuleConfig, error) {
	var config DingTalkRuleConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析钉钉规则配置失败: %v", err)
	}
	return &config, nil
}
