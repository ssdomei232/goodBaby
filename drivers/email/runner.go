package email

import (
	"context"

	"github.com/ssdomei232/goodBaby/model"
)

// EmailExecutor 邮件执行器
type EmailExecutor struct{}

func (e *EmailExecutor) GetType() string {
	return RuleType
}

func (e *EmailExecutor) Execute(ctx context.Context, task *model.RuleTask) error {
	account, err := ParseAccountConfig(task.AccountConfig())
	if err != nil {
		return err
	}

	config, err := ParseRuleConfig(task.Rule.ConfigJson)
	if err != nil {
		return err
	}

	return SendMail(ctx, account, config)
}
