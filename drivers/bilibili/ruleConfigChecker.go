package bilibili

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/internal/meta"
)

// BilibiliDynamicRuleValidator B站动态规则验证器
type BilibiliDynamicRuleValidator struct{}

func (v *BilibiliDynamicRuleValidator) GetType() string {
	return RuleTypeDynamic
}

func (v *BilibiliDynamicRuleValidator) Validate(configJSON string) error {
	var config BiliDynamicConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析B站动态规则配置失败: %v", err)
	}

	if config.Msg == "" {
		return fmt.Errorf("B站动态规则配置中 msg 不能为空")
	}

	if len([]rune(config.Msg)) > 1000 {
		return fmt.Errorf("B站动态内容过长")
	}

	return nil
}

func (v *BilibiliDynamicRuleValidator) Meta() meta.RuleMeta {
	return meta.RuleMeta{
		Type:        RuleTypeDynamic,
		Label:       "发送 B 站动态",
		Description: "触发时以关联的 B 站账号发送一条动态。",
		Docs:        "docs/bilibili-config.md",
		AccountType: AccountType,
		Fields: []meta.Field{
			{
				Key:         "msg",
				Label:       "动态内容",
				Type:        meta.FieldTextarea,
				Required:    true,
				Placeholder: "要发送的动态正文",
			},
		},
	}
}
