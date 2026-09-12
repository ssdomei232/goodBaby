package gateway

import "sync"

var (
	once     sync.Once
	registry *Registry
)

// InitGatewayRegistry 返回全局唯一的消息网关注册表
func InitGatewayRegistry() *Registry {
	once.Do(func() {
		r := NewRegistry()

		// 注册所有消息网关
		r.Register(&WebhookGateway{})
		// 未来添加新网关类型时，在这里注册即可

		registry = r
	})
	return registry
}
