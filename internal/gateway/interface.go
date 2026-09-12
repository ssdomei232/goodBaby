// Package gateway 定义消息网关的扩展接口与注册表。
//
// 与 drivers 的分工：driver 负责把一条规则执行出去(发邮件、发消息……)，
// gateway 负责把外部系统投递进来的消息派发给绑定在网关上的规则。
// 新增一种网关只要实现 Gateway 并在 reg.go 里注册。
package gateway

import (
	"context"
	"fmt"
	"sort"

	"github.com/ssdomei232/goodBaby/internal/meta"
	"github.com/ssdomei232/goodBaby/model"
)

// Message 是外部系统投递进来的一条消息
type Message struct {
	Title   string
	Content string
}

// Result 是一次投递的汇总结果
type Result struct {
	Total  int      `json:"total"`
	Failed []string `json:"failed"`
}

// Task 是一次投递所需的全部输入：网关与待触发的规则都由上层从数据库取出
type Task struct {
	Gateway *model.MessageGateway
	Rules   []model.GatewayRule
	Message Message
}

// Gateway 一种消息网关
type Gateway interface {
	// GetType 获取网关类型标识
	GetType() string
	// Meta 返回给 WebUI 展示用的元数据
	Meta() meta.GatewayMeta
	// Deliver 把投递进来的消息派发给 task 里的规则
	Deliver(ctx context.Context, task *Task) (*Result, error)
}

// Registry 网关注册表
type Registry struct {
	gateways map[string]Gateway
}

// NewRegistry 创建新的网关注册表
func NewRegistry() *Registry {
	return &Registry{
		gateways: make(map[string]Gateway),
	}
}

// Register 注册消息网关
func (r *Registry) Register(g Gateway) {
	r.gateways[g.GetType()] = g
}

// Resolve 按类型取网关实现，类型为空时回退到默认类型
func (r *Registry) Resolve(gatewayType string) (Gateway, bool) {
	if gatewayType == "" {
		gatewayType = model.GatewayTypeWebhook
	}

	g, exists := r.gateways[gatewayType]
	return g, exists
}

// Deliver 按网关类型派发一次投递
func (r *Registry) Deliver(ctx context.Context, gatewayType string, task *Task) (*Result, error) {
	g, exists := r.Resolve(gatewayType)
	if !exists {
		return nil, fmt.Errorf("不支持的消息网关类型: %s", gatewayType)
	}

	return g.Deliver(ctx, task)
}

// GetSupportedTypes 获取所有支持的网关类型
func (r *Registry) GetSupportedTypes() []string {
	types := make([]string, 0, len(r.gateways))
	for t := range r.gateways {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}

// Metas 返回所有网关类型的元数据，按类型名排序
func (r *Registry) Metas() []meta.GatewayMeta {
	metas := make([]meta.GatewayMeta, 0, len(r.gateways))
	for _, g := range r.gateways {
		metas = append(metas, g.Meta())
	}
	sort.Slice(metas, func(i, j int) bool { return metas[i].Type < metas[j].Type })
	return metas
}
