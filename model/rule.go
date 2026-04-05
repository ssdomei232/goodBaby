package model

type Rule struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	UID        uint   `json:"uid"`
	TimerID    uint   `json:"timer_id"`
	Name       string `json:"name"`
	AccountID  uint   `json:"account_id"`
	Type       string `json:"type"`
	ConfigJson string `json:"config_json"`
}
