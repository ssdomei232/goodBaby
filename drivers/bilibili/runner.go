package bilibili

import (
	"log"

	"github.com/ssdomei232/goodBaby/model"
)

// BilibiliDynamicExecutor B站动态执行器
type BilibiliDynamicExecutor struct{}

func (e *BilibiliDynamicExecutor) GetType() string {
	return "bilibili-dynamic"
}

func (e *BilibiliDynamicExecutor) Execute(rule *model.Rule) error {
	log.Printf("执行B站动态规则: %s (ID: %d)", rule.Name, rule.ID)

	SendBiliDynamicMsg(rule)

	return nil
}
