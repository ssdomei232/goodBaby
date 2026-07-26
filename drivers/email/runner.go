package email

import (
	"context"
	"log"

	"github.com/ssdomei232/goodBaby/model"
)

// EmailExecutor 邮件执行器
type EmailExecutor struct{}

func (e *EmailExecutor) GetType() string {
	return RuleType
}

func (e *EmailExecutor) Execute(ctx context.Context, rule *model.Rule) error {
	log.Printf("执行邮件规则: %s (ID: %d)", rule.Name, rule.ID)

	return SendMail(ctx, rule)
}
