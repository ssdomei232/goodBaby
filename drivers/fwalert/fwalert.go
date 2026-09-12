package fwalert

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/ssdomei232/goodBaby/internal/retry"
)

// Send 发送一条饭碗警告
func Send(ctx context.Context, config *FwAlertRuleConfig) error {
	if err := retry.Do(ctx, func() error {
		return send(config)
	}); err != nil {
		return fmt.Errorf("发送饭碗警告失败: %w", err)
	}
	return nil
}

// send 调用一次 Webhook
func send(config *FwAlertRuleConfig) error {
	resp, err := resty.New().R().
		SetBody(FwAlertRequest{Message: config.Msg}).
		Execute(resty.MethodPost, config.WebhookURL)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("Webhook 返回状态码 %s", resp.Status())
	}
	return nil
}
