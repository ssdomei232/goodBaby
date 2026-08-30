package rainyun

const (
	// AccountType Rainyun 账号类型标识
	AccountType = "rainyun"
	// RuleTypeWorkOrder 发送 Rainyun 工单消息的规则类型标识
	RuleTypeWorkOrder = "rainyun-workorder"
	RuleTypeRunAway   = "rainyun-runaway"
)

type RainyunWorkOrderRule struct {
	Title string `json:"title"`
	Msg   string `json:"msg"`
}

type RainyunAccount struct {
	APIKey string `json:"api_key"`
}

type RainyunRunAwayRule struct {
	IK string `json:"ik"`
}
