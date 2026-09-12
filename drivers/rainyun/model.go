package rainyun

const (
	// AccountType Rainyun 账号类型标识
	AccountType = "rainyun"
	// RuleTypeWorkOrder 发送 Rainyun 工单消息的规则类型标识
	RuleTypeWorkOrder = "rainyun-workorder"
	RuleTypeRunAway   = "rainyun-runaway"
)

// pageSize 拉取云服务器列表时每页的条数
const pageSize = 20

// runAwayConfirmText 跑路规则要求用户输入的免责声明
const runAwayConfirmText = "我已知晓"

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
