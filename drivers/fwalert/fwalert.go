package fwalert

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// send 发送消息
func send(config *FwAlertRuleConfig) error {
	restyClient := resty.New()
	req := restyClient.R()
	resp, err := req.Execute("GET", config.WebhookURL+"?message="+config.Msg)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("发送饭碗警告失败: %s", resp.Status())
	}
	return nil
}
