package accountConfigChecker

import (
	"github.com/ssdomei232/goodBaby/drivers/bilibili"
	"github.com/ssdomei232/goodBaby/drivers/email"
	"github.com/ssdomei232/goodBaby/drivers/github"
	"github.com/ssdomei232/goodBaby/drivers/onebot"
)

func InitValidatorRegistry() *ValidatorRegistry {
	registry := NewValidatorRegistry()

	// 注册所有账号验证器
	registry.Register(&bilibili.BilibiliAccountConfigValidator{})
	registry.Register(&email.EmailAccountConfigValidator{})
	registry.Register(&github.GitHubAccountConfigValidator{})
	registry.Register(&onebot.OneBotAccountConfigValidator{})

	return registry
}
