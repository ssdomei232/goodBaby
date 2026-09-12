package alidns

import (
	"context"

	"github.com/ssdomei232/goodBaby/model"
)

// DeleteAliDNSRecordExecutor 删除阿里云DNS记录执行器
type DeleteAliDNSRecordExecutor struct{}

func (e *DeleteAliDNSRecordExecutor) GetType() string {
	return RuleTypeDeleteRecord
}

func (e *DeleteAliDNSRecordExecutor) Execute(ctx context.Context, task *model.RuleTask) error {
	account, err := ParseAccountConfig(task.AccountConfig())
	if err != nil {
		return err
	}

	config, err := ParseDeleteRecordConfig(task.Rule.ConfigJson)
	if err != nil {
		return err
	}

	return DeleteRecord(ctx, account, config)
}
