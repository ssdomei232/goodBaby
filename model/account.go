package model

// Account 保存执行规则时需要用到的第三方凭据
type Account struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UID      uint   `gorm:"index" json:"uid"`
	Name     string `json:"name"`   // 账号名称
	Type     string `json:"type"`   // 账号类型，如 "bilibili"
	Config   string `json:"config"` // 存储账号相关配置，如 cookie 等
	CreateAt int64  `json:"create_at"`
}

// AccountRequest 创建/编辑账号的请求体
type AccountRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config string `json:"config"`
}

// Validate 校验账号请求中与类型无关的通用部分
func (r *AccountRequest) Validate() error {
	if r.Name == "" {
		return ErrValidation("账号名称不能为空")
	}
	if len(r.Name) > 64 {
		return ErrValidation("账号名称过长")
	}
	if r.Type == "" {
		return ErrValidation("账号类型不能为空")
	}
	if r.Config == "" {
		return ErrValidation("账号配置不能为空")
	}
	return nil
}
