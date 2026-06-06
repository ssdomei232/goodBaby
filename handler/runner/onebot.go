package runner

import (
	"log"

	"github.com/ssdomei232/goodBaby/drivers/onebot"
	"github.com/ssdomei232/goodBaby/model"
)

// OneBotExecutor OneBot执行器
type OneBotExecutor struct{}

func (e *OneBotExecutor) GetType() string {
	return "onebot"
}

func (e *OneBotExecutor) Execute(rule *model.Rule) error {
	log.Printf("执行OneBot规则: %s (ID: %d)", rule.Name, rule.ID)

	onebot.SendOneBotMsg(rule)

	return nil
}
