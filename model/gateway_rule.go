package model

// GatewayRule 是挂在消息网关上的通知规则。
//
// 与 Rule 的区别只有触发来源：Rule 由定时器到期触发，
// GatewayRule 由外部系统投递到网关的消息触发，两者互不干扰。
type GatewayRule struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	UID        uint   `gorm:"index" json:"uid"`
	GatewayID  uint   `gorm:"index" json:"gateway_id"`
	Name       string `json:"name"`
	AccountID  uint   `json:"account_id"`
	Type       string `json:"type"`
	ConfigJson string `json:"config_json"`
	Enabled    bool   `gorm:"default:true" json:"enabled"`
	CreateAt   int64  `json:"create_at"`
}

// AsRule 把网关规则换成统一的执行结构。
//
// 执行层只认识 Rule：投递消息前把外部消息覆盖到 ConfigJson 上，
// 再交给同一个 runner 执行即可。
func (r *GatewayRule) AsRule() *Rule {
	return &Rule{
		UID:        r.UID,
		Name:       r.Name,
		AccountID:  r.AccountID,
		Type:       r.Type,
		ConfigJson: r.ConfigJson,
		Enabled:    r.Enabled,
		CreateAt:   r.CreateAt,
	}
}

// GatewayRuleRequest 创建/编辑消息网关规则的请求体
type GatewayRuleRequest struct {
	Name       string `json:"name"`
	GatewayID  uint   `json:"gateway_id"`
	AccountID  uint   `json:"account_id"`
	Type       string `json:"type"`
	ConfigJson string `json:"config_json"`
	Enabled    *bool  `json:"enabled"`
}

// Validate 校验网关规则请求中与类型无关的通用部分
func (r *GatewayRuleRequest) Validate() error {
	if r.Name == "" {
		return ErrValidation("规则名称不能为空")
	}
	if len(r.Name) > 64 {
		return ErrValidation("规则名称过长")
	}
	if r.Type == "" {
		return ErrValidation("规则类型不能为空")
	}
	if r.GatewayID == 0 {
		return ErrValidation("规则必须关联一个消息网关")
	}
	if r.ConfigJson == "" {
		return ErrValidation("规则配置不能为空")
	}
	return nil
}
