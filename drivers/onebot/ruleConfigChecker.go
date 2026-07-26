package onebot

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/internal/meta"
)

// OneBotRuleValidator OneBot规则验证器
type OneBotRuleValidator struct{}

func (v *OneBotRuleValidator) GetType() string {
	return RuleType
}

func (v *OneBotRuleValidator) Validate(configJSON string) error {
	var config OneBotConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析OneBot规则配置失败: %v", err)
	}

	if config.Msg == "" {
		return fmt.Errorf("OneBot规则配置中 msg 不能为空")
	}

	if len(config.SendGroups) == 0 && len(config.SendUsers) == 0 {
		return fmt.Errorf("OneBot规则配置中 send_groups 与 send_users 不能同时为空")
	}

	for _, id := range append(append([]int64{}, config.SendGroups...), config.SendUsers...) {
		if id <= 0 {
			return fmt.Errorf("OneBot规则配置中存在非法的群号/QQ号: %d", id)
		}
	}

	return nil
}

func (v *OneBotRuleValidator) Meta() meta.RuleMeta {
	return meta.RuleMeta{
		Type:        RuleType,
		Label:       "发送 QQ 消息",
		Description: "触发时通过 OneBot 向指定的群或好友发送消息。",
		AccountType: AccountType,
		Fields: []meta.Field{
			{Key: "msg", Label: "消息内容", Type: meta.FieldTextarea, Required: true},
			{Key: "send_groups", Label: "群号", Type: meta.FieldNumberList, Placeholder: "123456789"},
			{Key: "send_users", Label: "好友 QQ 号", Type: meta.FieldNumberList, Placeholder: "123456789"},
		},
	}
}
