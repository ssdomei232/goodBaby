package rainyun

import (
	"context"
	"log"

	"github.com/ssdomei232/goodBaby/model"
)

// RainyunExecutor 雨云执行器
type RainyunWorkorderExecutor struct{}

func (e *RainyunWorkorderExecutor) GetType() string {
	return RuleType
}

func (e *RainyunWorkorderExecutor) Execute(ctx context.Context, rule *model.Rule) error {
	log.Printf("执行雨云工单规则: %s", rule.Name)

	return sendRainyunWorkorderMsg(ctx, rule)
}
