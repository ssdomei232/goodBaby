package model

// Rule 描述一个 Timer 到期后要执行的动作
type Rule struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	UID       uint   `gorm:"index" json:"uid"`
	TimerID   uint   `gorm:"index" json:"timer_id"`
	GatewayID uint   `gorm:"index" json:"gateway_id"`
	Name      string `json:"name"`
	// 关联账号，部分规则类型(如钉钉)不需要账号，此时为 0
	AccountID uint `json:"account_id"`
	// 规则类型，见 internal/ruleConfigChecker 注册表
	Type string `json:"type"`
	// 规则配置，JSON 字符串
	ConfigJson string `json:"config_json"`
	Enabled    bool   `gorm:"default:true" json:"enabled"`
	CreateAt   int64  `json:"create_at"`
}

// RuleRequest 创建/编辑规则的请求体
type RuleRequest struct {
	Name       string `json:"name"`
	TimerID    uint   `json:"timer_id"`
	GatewayID  uint   `json:"gateway_id"`
	AccountID  uint   `json:"account_id"`
	Type       string `json:"type"`
	ConfigJson string `json:"config_json"`
	Enabled    *bool  `json:"enabled"`
}

// Validate 校验规则请求中与类型无关的通用部分
func (r *RuleRequest) Validate() error {
	if r.Name == "" {
		return ErrValidation("规则名称不能为空")
	}
	if len(r.Name) > 64 {
		return ErrValidation("规则名称过长")
	}
	if r.Type == "" {
		return ErrValidation("规则类型不能为空")
	}
	if r.TimerID == 0 && r.GatewayID == 0 {
		return ErrValidation("必须关联一个定时器或消息网关")
	}
	if r.TimerID != 0 && r.GatewayID != 0 {
		return ErrValidation("定时器和消息网关只能选择一个")
	}
	if r.ConfigJson == "" {
		return ErrValidation("规则配置不能为空")
	}
	return nil
}
