package fwalert

import (
	"context"
	"log"

	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

// FwalertExecutor 饭碗警告执行器
type FwalertExecutor struct{}

func (e *FwalertExecutor) GetType() string {
	return RuleTypeFwalert
}

func (e *FwalertExecutor) Execute(ctx context.Context, rule *model.Rule) error {
	log.Printf("执行饭碗警告规则: %s (ID: %d)", rule.Name, rule.ID)

	config, err := ParseRuleConfig(rule.ConfigJson)
	if err != nil {
		return err
	}

	return retry.Do(ctx, func() error {
		return send(config)
	})
}
