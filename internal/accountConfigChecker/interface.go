package accountConfigChecker

import (
	"encoding/json"
	"fmt"
)

// AccountValidator 账号验证器接口
type AccountValidator interface {
	// Validate 验证账号配置是否合法
	Validate(account string) error
	// GetType 获取验证器支持的账号类型
	GetType() string
}

// ValidatorRegistry 验证器注册表
type ValidatorRegistry struct {
	validators map[string]AccountValidator
}

// NewValidatorRegistry 创建新的验证器注册表
func NewValidatorRegistry() *ValidatorRegistry {
	return &ValidatorRegistry{
		validators: make(map[string]AccountValidator),
	}
}

// Register 注册账号验证器
func (vr *ValidatorRegistry) Register(validator AccountValidator) {
	vr.validators[validator.GetType()] = validator
}

// Validate 根据账号类型验证配置
func (vr *ValidatorRegistry) Validate(accountType string, config string) error {
	validator, exists := vr.validators[accountType]
	if !exists {
		return fmt.Errorf("不支持的账号类型: %s", accountType)
	}

	if config == "" {
		return fmt.Errorf("账号配置不能为空")
	}

	// 首先验证 JSON 格式是否正确
	var js json.RawMessage
	if err := json.Unmarshal([]byte(config), &js); err != nil {
		return fmt.Errorf("账号配置 JSON 格式错误: %v", err)
	}

	// 调用具体验证器的验证逻辑
	return validator.Validate(config)
}

// GetSupportedTypes 获取所有支持的账号类型
func (vr *ValidatorRegistry) GetSupportedTypes() []string {
	types := make([]string, 0, len(vr.validators))
	for t := range vr.validators {
		types = append(types, t)
	}
	return types
}
