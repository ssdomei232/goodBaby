package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/handler/runner"
	"github.com/ssdomei232/goodBaby/internal/meta"
	"github.com/ssdomei232/goodBaby/model"
)

// messageFields 规则配置里可以被投递消息覆盖的字段名
var messageFields = []string{"msg", "message", "body"}

// titleField 规则配置里可以被投递标题覆盖的字段名
const titleField = "title"

// TestMessage 是「测试」网关规则时使用的示例投递内容，
// 让用户在页面上就能看到消息被替换后的效果。
var TestMessage = Message{
	Title:   "goodBaby 测试标题",
	Content: "这是一条来自消息网关的测试消息",
}

// WebhookGateway 默认网关：外部系统 POST 一条消息来触发规则
type WebhookGateway struct{}

func (g *WebhookGateway) GetType() string {
	return model.GatewayTypeWebhook
}

func (g *WebhookGateway) Meta() meta.GatewayMeta {
	return meta.GatewayMeta{
		Type:        model.GatewayTypeWebhook,
		Label:       "Webhook",
		Description: "为外部系统生成一个 Webhook 地址，POST 一条消息即可触发绑定在该网关上的规则。",
		Docs:        "docs/gateway-config.md",
		PayloadHint: `{"title": "监控告警", "message": "服务已恢复"}`,
	}
}

// Deliver 逐条执行绑定在网关上的规则
//
// 单条规则失败不影响其它规则，失败原因汇总在 Result 里返回。
func (g *WebhookGateway) Deliver(ctx context.Context, task *Task) (*Result, error) {
	// 初始化为空切片：nil 切片会被序列化成 JSON null，前端读 failed.length 时会报错
	failed := []string{}

	for i := range task.Rules {
		rule := task.Rules[i]

		executed := rule
		configJSON, applied, err := ApplyMessage(rule.ConfigJson, task.Message)
		if err != nil {
			failed = append(failed, rule.Name+": "+err.Error())
			continue
		}
		if applied {
			// 用投递进来的消息执行，但不改动数据库里的规则配置
			executed.ConfigJson = configJSON
		}
		if err := runner.ExecuteGatewayRuleWithContext(ctx, &executed, model.TriggerGateway); err != nil {
			failed = append(failed, rule.Name+": "+err.Error())
		}
	}

	return &Result{Total: len(task.Rules), Failed: failed}, nil
}

// ApplyMessage 把投递进来的消息覆盖到规则配置的消息字段上。
//
// 返回覆盖后的配置；applied 为 false 表示配置里没有可覆盖的消息字段
// (例如「公开 GitHub 仓库」「删除 DNS 记录」)，此时按规则里保存的配置执行。
func ApplyMessage(configJSON string, msg Message) (string, bool, error) {
	var config map[string]any
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return configJSON, false, fmt.Errorf("解析规则配置失败: %v", err)
	}

	matched := false
	for _, key := range messageFields {
		if _, ok := config[key]; ok {
			if msg.Content != "" {
				config[key] = msg.Content
			}
			matched = true
		}
	}

	if msg.Title != "" {
		if _, ok := config[titleField]; ok {
			config[titleField] = msg.Title
			matched = true
		}
	}

	if !matched {
		return configJSON, false, nil
	}

	out, err := json.Marshal(config)
	if err != nil {
		return configJSON, false, fmt.Errorf("生成规则配置失败: %v", err)
	}
	return string(out), true, nil
}
