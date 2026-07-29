package dingtalk

import (
	"context"
	"log"

	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

// DingTalkExecutor 钉钉机器人执行器
type DingTalkExecutor struct{}

func (e *DingTalkExecutor) GetType() string {
	return RuleType
}

func (e *DingTalkExecutor) Execute(ctx context.Context, rule *model.Rule) error {
	log.Printf("执行钉钉规则: %s (ID: %d)", rule.Name, rule.ID)

	config, err := ParseRuleConfig(rule.ConfigJson)
	if err != nil {
		return err
	}

	return retry.Do(ctx, func() error {
		return send(&config.DingTalkConfig, config.Title, config.Msg)
	})
}
