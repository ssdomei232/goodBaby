package model

// RuleTask 是一次规则执行所需的全部输入。
//
// 规则本体与关联账号都由上层(handler/runner、消息网关)从数据库取出后传进来，
// driver 只负责把动作做出去，所以底层实现里不会再出现任何数据库操作。
type RuleTask struct {
	// Rule 规则本体，永远不会为空
	Rule *Rule
	// Account 规则关联的账号；该规则类型不需要账号时为 nil
	Account *Account
}

// AccountConfig 返回关联账号的配置 JSON，没有关联账号时返回空字符串
func (t *RuleTask) AccountConfig() string {
	if t == nil || t.Account == nil {
		return ""
	}
	return t.Account.Config
}
