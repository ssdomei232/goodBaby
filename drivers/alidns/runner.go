package alidns

import (
	"context"
	"log"

	"github.com/ssdomei232/goodBaby/model"
)

// DeleteAliDNSRecordExecutor 删除阿里云DNS记录执行器
type DeleteAliDNSRecordExecutor struct{}

func (e *DeleteAliDNSRecordExecutor) GetType() string {
	return RuleTypeDeleteRecord
}

func (e *DeleteAliDNSRecordExecutor) Execute(ctx context.Context, rule *model.Rule) error {
	log.Printf("执行删除阿里云DNS记录规则: %s (ID: %d)", rule.Name, rule.ID)

	return deleteAliDNSRecord(ctx, rule)
}
