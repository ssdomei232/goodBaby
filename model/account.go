package model

type Account struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UID    uint   `json:"uid"`
	Name   string `json:"name"`   // 账号名称
	Type   string `json:"type"`   // 账号类型，如 "bilibili"
	Config string `json:"config"` // 存储账号相关配置，如cookie等
}
