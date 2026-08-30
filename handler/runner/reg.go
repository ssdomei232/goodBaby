package runner

import (
	"github.com/ssdomei232/goodBaby/drivers/alidns"
	"github.com/ssdomei232/goodBaby/drivers/bilibili"
	"github.com/ssdomei232/goodBaby/drivers/dingtalk"
	"github.com/ssdomei232/goodBaby/drivers/email"
	"github.com/ssdomei232/goodBaby/drivers/fwalert"
	"github.com/ssdomei232/goodBaby/drivers/github"
	"github.com/ssdomei232/goodBaby/drivers/onebot"
	"github.com/ssdomei232/goodBaby/drivers/rainyun"
)

// InitExecutorRegistry 初始化执行器注册表并注册所有执行器
func InitExecutorRegistry() *ExecutorRegistry {
	registryOnce.Do(func() {
		registry := NewExecutorRegistry()

		// 注册所有规则执行器
		registry.Register(&bilibili.BilibiliDynamicExecutor{})
		registry.Register(&bilibili.BilibiliPrivateMessageExecutor{})
		registry.Register(&email.EmailExecutor{})
		registry.Register(&github.GithubMakeRepoPublicExecutor{})
		registry.Register(&onebot.OneBotExecutor{})
		registry.Register(&dingtalk.DingTalkExecutor{})
		registry.Register(&fwalert.FwalertExecutor{})
		registry.Register(&rainyun.RainyunWorkorderExecutor{})
		registry.Register(&rainyun.RainyunRunAwayExecutor{})
		registry.Register(&alidns.DeleteAliDNSRecordExecutor{})
		// 未来添加新规则类型时，在这里注册即可

		globalExecutorRegistry = registry
	})
	return globalExecutorRegistry
}
