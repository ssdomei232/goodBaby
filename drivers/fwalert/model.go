package fwalert

// RuleType 饭碗警告规则类型标识
const RuleType = "fwalert"

type FwAlertRuleConfig struct {
	WebhookURL string `json:"webhook_url"`
	Msg        string `json:"msg"`
}

type FwAlertRequest struct {
	Message string `json:"message"`
}
