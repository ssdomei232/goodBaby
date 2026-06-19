package onebot

import (
	"encoding/json"
	"fmt"
)

// OneBotRuleValidator OneBot规则验证器
type OneBotRuleValidator struct{}

func (v *OneBotRuleValidator) GetType() string {
	return "onebot"
}

func (v *OneBotRuleValidator) Validate(configJSON string) error {
	var config OneBotConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析OneBot规则配置失败: %v", err)
	}

	if config.Msg == "" {
		return fmt.Errorf("OneBot规则配置中 msg 不能为空")
	}

	if len(config.SendGroups) == 0 {
		return fmt.Errorf("OneBot规则配置中 send_groups 不能为空")
	}

	return nil
}
