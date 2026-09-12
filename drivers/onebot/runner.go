package onebot

import (
	"context"

	"github.com/ssdomei232/goodBaby/model"
)

// OneBotExecutor OneBot执行器
type OneBotExecutor struct{}

func (e *OneBotExecutor) GetType() string {
	return RuleType
}

func (e *OneBotExecutor) Execute(ctx context.Context, task *model.RuleTask) error {
	account, err := ParseAccountConfig(task.AccountConfig())
	if err != nil {
		return err
	}

	config, err := ParseRuleConfig(task.Rule.ConfigJson)
	if err != nil {
		return err
	}

	return SendMessage(ctx, account, config)
}
