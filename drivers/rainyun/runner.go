package rainyun

import (
	"context"
	"fmt"

	"github.com/ssdomei232/goodBaby/model"
)

// RainyunWorkorderExecutor 雨云工单执行器
type RainyunWorkorderExecutor struct{}

func (e *RainyunWorkorderExecutor) GetType() string {
	return RuleTypeWorkOrder
}

func (e *RainyunWorkorderExecutor) Execute(ctx context.Context, task *model.RuleTask) error {
	account, err := ParseAccountConfig(task.AccountConfig())
	if err != nil {
		return err
	}

	config, err := ParseWorkOrderConfig(task.Rule.ConfigJson)
	if err != nil {
		return err
	}

	return SendWorkOrder(ctx, account, config)
}

// RainyunRunAwayExecutor 雨云跑路执行器
type RainyunRunAwayExecutor struct{}

func (e *RainyunRunAwayExecutor) GetType() string {
	return RuleTypeRunAway
}

func (e *RainyunRunAwayExecutor) Execute(ctx context.Context, task *model.RuleTask) error {
	account, err := ParseAccountConfig(task.AccountConfig())
	if err != nil {
		return err
	}

	// 跑路规则只要求用户填一句免责声明，执行时重新校验一次
	config, err := ParseRunAwayConfig(task.Rule.ConfigJson)
	if err != nil {
		return err
	}
	if config.IK != runAwayConfirmText {
		return fmt.Errorf("请输入正确内容")
	}

	return RunAway(ctx, account)
}
