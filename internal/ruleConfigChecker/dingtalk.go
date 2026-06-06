package ruleConfigChecker

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/drivers/dingtalk"
)

// DingTalkRuleValidator 钉钉规则验证器
type DingTalkRuleValidator struct{}

func (v *DingTalkRuleValidator) GetType() string {
	return "dingtalk"
}

func (v *DingTalkRuleValidator) Validate(configJSON string) error {
	var config dingtalk.DingTalkConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析钉钉规则配置失败: %v", err)
	}

	if config.AccessToken == "" {
		return fmt.Errorf("钉钉规则配置中 access_token 不能为空")
	}

	return nil
}
