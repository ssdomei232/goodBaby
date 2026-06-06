package runner

import (
	"log"

	"github.com/ssdomei232/goodBaby/drivers/email"
	"github.com/ssdomei232/goodBaby/model"
)

// EmailExecutor 邮件执行器
type EmailExecutor struct{}

func (e *EmailExecutor) GetType() string {
	return "email"
}

func (e *EmailExecutor) Execute(rule *model.Rule) error {
	log.Printf("执行邮件规则: %s (ID: %d)", rule.Name, rule.ID)

	email.SendMail(rule)

	return nil
}
