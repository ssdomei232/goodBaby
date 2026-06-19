package bilibili

import (
	"encoding/json"
	"fmt"
)

// BilibiliDynamicRuleValidator B站动态规则验证器
type BilibiliDynamicRuleValidator struct{}

func (v *BilibiliDynamicRuleValidator) GetType() string {
	return "bilibili-dynamic"
}

func (v *BilibiliDynamicRuleValidator) Validate(configJSON string) error {
	var config BiliDynamicConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析B站动态规则配置失败: %v", err)
	}

	if config.Msg == "" {
		return fmt.Errorf("B站动态规则配置中 msg 不能为空")
	}

	return nil
}
