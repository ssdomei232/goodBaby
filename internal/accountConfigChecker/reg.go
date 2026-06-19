package accountConfigChecker

func InitValidatorRegistry() *ValidatorRegistry {
	registry := NewValidatorRegistry()

	// 注册所有规则验证器
	// registry.Register(&bilibili.BilibiliAccountValidator{})

	return registry
}
