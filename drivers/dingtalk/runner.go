package dingtalk

import (
	"context"

	"github.com/ssdomei232/goodBaby/model"
)

// DingTalkExecutor 钉钉机器人执行器
type DingTalkExecutor struct{}

func (e *DingTalkExecutor) GetType() string {
	return RuleType
}

func (e *DingTalkExecutor) Execute(ctx context.Context, task *model.RuleTask) error {
	// 钉钉机器人凭据直接写在规则配置里，不需要关联账号
	config, err := ParseRuleConfig(task.Rule.ConfigJson)
	if err != nil {
		return err
	}

	return SendMessage(ctx, config)
}
