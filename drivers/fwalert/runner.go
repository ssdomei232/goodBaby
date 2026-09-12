package fwalert

import (
	"context"

	"github.com/ssdomei232/goodBaby/model"
)

// FwalertExecutor 饭碗警告执行器
type FwalertExecutor struct{}

func (e *FwalertExecutor) GetType() string {
	return RuleType
}

func (e *FwalertExecutor) Execute(ctx context.Context, task *model.RuleTask) error {
	// 饭碗警告的凭据就是 Webhook 地址，直接写在规则配置里
	config, err := ParseRuleConfig(task.Rule.ConfigJson)
	if err != nil {
		return err
	}

	return Send(ctx, config)
}
