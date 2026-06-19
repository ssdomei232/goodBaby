package ruleConfigChecker

import (
	"encoding/json"
	"fmt"
)

// RuleValidator 规则验证器接口
type RuleValidator interface {
	// Validate 验证规则配置是否合法
	Validate(configJSON string) error
	// GetType 获取验证器支持的规则类型
	GetType() string
}

// ValidatorRegistry 验证器注册表
type ValidatorRegistry struct {
	validators map[string]RuleValidator
}

// NewValidatorRegistry 创建新的验证器注册表
func NewValidatorRegistry() *ValidatorRegistry {
	return &ValidatorRegistry{
		validators: make(map[string]RuleValidator),
	}
}

// Register 注册规则验证器
func (vr *ValidatorRegistry) Register(validator RuleValidator) {
	vr.validators[validator.GetType()] = validator
}

// Validate 根据规则类型验证配置
func (vr *ValidatorRegistry) Validate(ruleType string, configJSON string) error {
	validator, exists := vr.validators[ruleType]
	if !exists {
		return fmt.Errorf("不支持的规则类型: %s", ruleType)
	}

	if configJSON == "" {
		return fmt.Errorf("规则配置不能为空")
	}

	// 首先验证 JSON 格式是否正确
	var js json.RawMessage
	if err := json.Unmarshal([]byte(configJSON), &js); err != nil {
		return fmt.Errorf("规则配置 JSON 格式错误: %v", err)
	}

	// 调用具体验证器的验证逻辑
	return validator.Validate(configJSON)
}

// GetSupportedTypes 获取所有支持的规则类型
func (vr *ValidatorRegistry) GetSupportedTypes() []string {
	types := make([]string, 0, len(vr.validators))
	for t := range vr.validators {
		types = append(types, t)
	}
	return types
}
