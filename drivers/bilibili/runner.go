package bilibili

import (
	"context"

	"github.com/ssdomei232/goodBaby/model"
)

// BilibiliDynamicExecutor B站动态执行器
type BilibiliDynamicExecutor struct{}

func (e *BilibiliDynamicExecutor) GetType() string {
	return RuleTypeDynamic
}

func (e *BilibiliDynamicExecutor) Execute(ctx context.Context, task *model.RuleTask) error {
	account, err := ParseAccountConfig(task.AccountConfig())
	if err != nil {
		return err
	}

	config, err := ParseDynamicConfig(task.Rule.ConfigJson)
	if err != nil {
		return err
	}

	return SendDynamic(ctx, account, config)
}

// BilibiliPrivateMessageExecutor B站私信执行器
type BilibiliPrivateMessageExecutor struct{}

func (e *BilibiliPrivateMessageExecutor) GetType() string {
	return RuleTypePrivateMessage
}

func (e *BilibiliPrivateMessageExecutor) Execute(ctx context.Context, task *model.RuleTask) error {
	account, err := ParseAccountConfig(task.AccountConfig())
	if err != nil {
		return err
	}

	config, err := ParsePrivateMessageConfig(task.Rule.ConfigJson)
	if err != nil {
		return err
	}

	return SendPrivateMessage(ctx, account, config)
}
