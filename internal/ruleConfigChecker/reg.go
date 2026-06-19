package ruleConfigChecker

import (
	"github.com/ssdomei232/goodBaby/drivers/bilibili"
	"github.com/ssdomei232/goodBaby/drivers/dingtalk"
	"github.com/ssdomei232/goodBaby/drivers/email"
	"github.com/ssdomei232/goodBaby/drivers/github"
	"github.com/ssdomei232/goodBaby/drivers/onebot"
)

// InitValidatorRegistry 初始化验证器注册表并注册所有验证器
func InitValidatorRegistry() *ValidatorRegistry {
	registry := NewValidatorRegistry()

	// 注册所有规则验证器
	registry.Register(&bilibili.BilibiliDynamicRuleValidator{})
	registry.Register(&email.EmailRuleValidator{})
	registry.Register(&onebot.OneBotRuleValidator{})
	registry.Register(&dingtalk.DingTalkRuleValidator{})
	registry.Register(&github.GithubMakeRepositoryPublicRuleValidator{})

	return registry
}
