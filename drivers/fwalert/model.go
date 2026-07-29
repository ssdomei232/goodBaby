package fwalert

const RuleTypeFwalert = "fwalert"

type FwAlertRuleConfig struct {
	WebhookURL string `json:"webhook_url"`
	Msg        string `json:"msg"`
}
