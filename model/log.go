package model

// 执行日志的触发来源
const (
	TriggerTimer  = "timer"  // Timer 到期自动触发
	TriggerManual = "manual" // 用户在 WebUI 手动测试
	TriggerRemind = "remind" // 提醒通知
)

// ExecutionLog 记录每一次规则执行/提醒的结果，供 WebUI 展示
type ExecutionLog struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UID      uint   `gorm:"index" json:"uid"`
	RuleID   uint   `json:"rule_id"`
	RuleName string `json:"rule_name"`
	RuleType string `json:"rule_type"`
	TimerID  uint   `json:"timer_id"`
	Trigger  string `json:"trigger"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	CreateAt int64  `gorm:"index" json:"create_at"`
}
