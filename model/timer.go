package model

type Timer struct {
	ID                  uint   `gorm:"primaryKey" json:"id"`
	UID                 uint   `json:"uid"`
	Name                string `json:"name"`
	SignDerationSeconds int64  `json:"sign_deration_seconds"`
	LastSign            int64  `json:"last_sign"`
	RemindTimeSeconds   int64  `json:"remind_time_seconds"` // 提前多少秒提醒
}
