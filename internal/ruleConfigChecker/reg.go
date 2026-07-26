package ruleConfigChecker

import (
	"sync"

	"github.com/ssdomei232/goodBaby/drivers/bilibili"
	"github.com/ssdomei232/goodBaby/drivers/dingtalk"
	"github.com/ssdomei232/goodBaby/drivers/email"
	"github.com/ssdomei232/goodBaby/drivers/github"
	"github.com/ssdomei232/goodBaby/drivers/onebot"
)

var (
	once     sync.Once
	registry *ValidatorRegistry
)

// InitValidatorRegistry 返回全局唯一的规则验证器注册表
func InitValidatorRegistry() *ValidatorRegistry {
	once.Do(func() {
		r := NewValidatorRegistry()

		// 注册所有规则验证器
		r.Register(&bilibili.BilibiliDynamicRuleValidator{})
		r.Register(&email.EmailRuleValidator{})
		r.Register(&onebot.OneBotRuleValidator{})
		r.Register(&dingtalk.DingTalkRuleValidator{})
		r.Register(&github.GithubMakeRepositoryPublicRuleValidator{})

		registry = r
	})
	return registry
}
