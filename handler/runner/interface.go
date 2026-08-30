package runner

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/ssdomei232/goodBaby/model"
)

// RuleExecutor 规则执行器接口
type RuleExecutor interface {
	// Execute 执行规则，ctx 决定了重试的最长时间
	Execute(ctx context.Context, rule *model.Rule) error
	// GetType 获取执行器支持的规则类型
	GetType() string
}

// ExecutorRegistry 执行器注册表
type ExecutorRegistry struct {
	executors map[string]RuleExecutor
}

// NewExecutorRegistry 创建新的执行器注册表
func NewExecutorRegistry() *ExecutorRegistry {
	return &ExecutorRegistry{
		executors: make(map[string]RuleExecutor),
	}
}

// Register 注册规则执行器
func (er *ExecutorRegistry) Register(executor RuleExecutor) {
	er.executors[executor.GetType()] = executor
}

// Execute 根据规则类型执行规则
func (er *ExecutorRegistry) Execute(ctx context.Context, rule *model.Rule) error {
	executor, exists := er.executors[rule.Type]
	if !exists {
		return fmt.Errorf("不支持的规则类型: %s", rule.Type)
	}

	return executor.Execute(ctx, rule)
}

// GetSupportedTypes 获取所有支持的规则类型
func (er *ExecutorRegistry) GetSupportedTypes() []string {
	types := make([]string, 0, len(er.executors))
	for t := range er.executors {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

var (
	registryOnce           sync.Once
	globalExecutorRegistry *ExecutorRegistry
)

// GetGlobalExecutorRegistry 获取全局执行器注册表
func GetGlobalExecutorRegistry() *ExecutorRegistry {
	if globalExecutorRegistry == nil {
		panic("执行器注册表未初始化，请先调用 InitExecutorRegistry()")
	}
	return globalExecutorRegistry
}
