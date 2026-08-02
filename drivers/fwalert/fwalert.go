package fwalert

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// send 发送消息
func send(config *FwAlertRuleConfig) error {
	restyClient := resty.New()
	req := restyClient.R()

	req.SetBody(FwAlertRequest{
		Message: config.Msg,
	})

	resp, err := req.Execute(resty.MethodPost, config.WebhookURL)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("发送饭碗警告失败: %s", resp.Status())
	}
	return nil
}
