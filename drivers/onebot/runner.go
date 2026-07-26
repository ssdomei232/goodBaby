package onebot

import (
	"context"
	"log"

	"github.com/ssdomei232/goodBaby/model"
)

// OneBotExecutor OneBot执行器
type OneBotExecutor struct{}

func (e *OneBotExecutor) GetType() string {
	return RuleType
}

func (e *OneBotExecutor) Execute(ctx context.Context, rule *model.Rule) error {
	log.Printf("执行OneBot规则: %s (ID: %d)", rule.Name, rule.ID)

	return SendOneBotMsg(ctx, rule)
}
