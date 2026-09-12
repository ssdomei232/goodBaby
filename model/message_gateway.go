package model

// GatewayTypeWebhook 默认的网关类型：外部系统 POST 一条消息来触发规则
const GatewayTypeWebhook = "webhook"

type MessageGateway struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	UID  uint   `gorm:"index" json:"uid"`
	Name string `json:"name"`
	// Type 网关类型，见 internal/gateway 注册表
	Type     string `gorm:"default:webhook" json:"type"`
	Token    string `gorm:"uniqueIndex;size:67" json:"token"`
	CreateAt int64  `json:"create_at"`
}

type MessageGatewayRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Validate 校验网关请求中与类型无关的通用部分
func (r *MessageGatewayRequest) Validate() error {
	if r.Name == "" {
		return ErrValidation("网关名称不能为空")
	}
	if len(r.Name) > 64 {
		return ErrValidation("网关名称过长")
	}
	return nil
}
