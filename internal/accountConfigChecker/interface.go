package accountConfigChecker

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/ssdomei232/goodBaby/internal/meta"
)

// AccountValidator 账号验证器接口
type AccountValidator interface {
	// Validate 验证账号配置是否合法
	Validate(config string) error
	// GetType 获取验证器支持的账号类型
	GetType() string
	// Meta 返回给 WebUI 渲染表单用的元数据
	Meta() meta.AccountMeta
}

// AccountTester 可选接口：支持真实连通性测试的账号类型实现它
type AccountTester interface {
	// Test 使用给定配置访问一次第三方服务，验证凭据是否可用
	Test(config string) error
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

// Test 对账号配置做一次真实的连通性测试
func (vr *ValidatorRegistry) Test(accountType string, config string) error {
	validator, exists := vr.validators[accountType]
	if !exists {
		return fmt.Errorf("不支持的账号类型: %s", accountType)
	}

	if err := vr.Validate(accountType, config); err != nil {
		return err
	}

	tester, ok := validator.(AccountTester)
	if !ok {
		return fmt.Errorf("账号类型 %s 暂不支持连通性测试", accountType)
	}
	return tester.Test(config)
}

// GetSupportedTypes 获取所有支持的账号类型
func (vr *ValidatorRegistry) GetSupportedTypes() []string {
	types := make([]string, 0, len(vr.validators))
	for t := range vr.validators {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

// Metas 返回所有账号类型的元数据，按类型名排序
func (vr *ValidatorRegistry) Metas() []meta.AccountMeta {
	metas := make([]meta.AccountMeta, 0, len(vr.validators))
	for _, v := range vr.validators {
		metas = append(metas, withTestable(v))
	}
	sort.Slice(metas, func(i, j int) bool { return metas[i].Type < metas[j].Type })
	return metas
}

// MetaOf 返回指定账号类型的元数据
func (vr *ValidatorRegistry) MetaOf(accountType string) (meta.AccountMeta, bool) {
	v, ok := vr.validators[accountType]
	if !ok {
		return meta.AccountMeta{}, false
	}
	return withTestable(v), true
}

// FieldsOf 返回指定账号类型的字段描述，类型不存在时返回 nil
func (vr *ValidatorRegistry) FieldsOf(accountType string) []meta.Field {
	m, ok := vr.MetaOf(accountType)
	if !ok {
		return nil
	}
	return m.Fields
}

func withTestable(v AccountValidator) meta.AccountMeta {
	m := v.Meta()
	if _, ok := v.(AccountTester); ok {
		m.Testable = true
	}
	return m
}
