package rainyun

const (
	// AccountType Rainyun 账号类型标识
	AccountType = "rainyun"
	// RuleType 发送 Rainyun 消息的规则类型标识
	RuleType = "rainyun-workorder"
)

type RainyunWorkOrderRule struct {
	Title string `json:"title"`
	Msg   string `json:"msg"`
}

type RainyunAccount struct {
	APIKey string `json:"api_key"`
}
