package accountConfigChecker

import (
	"sync"

	"github.com/ssdomei232/goodBaby/drivers/bilibili"
	"github.com/ssdomei232/goodBaby/drivers/email"
	"github.com/ssdomei232/goodBaby/drivers/github"
	"github.com/ssdomei232/goodBaby/drivers/onebot"
	"github.com/ssdomei232/goodBaby/drivers/rainyun"
)

var (
	once     sync.Once
	registry *ValidatorRegistry
)

// InitValidatorRegistry 返回全局唯一的账号验证器注册表
func InitValidatorRegistry() *ValidatorRegistry {
	once.Do(func() {
		r := NewValidatorRegistry()

		// 注册所有账号验证器
		r.Register(&bilibili.BilibiliAccountConfigValidator{})
		r.Register(&email.EmailAccountConfigValidator{})
		r.Register(&github.GitHubAccountConfigValidator{})
		r.Register(&onebot.OneBotAccountConfigValidator{})
		r.Register(&rainyun.RainyunAccountConfigValidator{})

		registry = r
	})
	return registry
}
