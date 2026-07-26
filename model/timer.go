package model

// Timer 是摇篮系统的核心：用户需要在 SignDerationSeconds 内完成一次签到，
// 否则挂在这个 Timer 下的所有规则都会被触发。
type Timer struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	UID         uint   `gorm:"index" json:"uid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// 是否启用，停用后不再提醒也不会触发
	Enabled bool `gorm:"default:true" json:"enabled"`
	// 签到有效期(秒)，超过这个时间没有签到就会触发规则
	SignDerationSeconds int64 `json:"sign_deration_seconds"`
	// 上次签到时间(Unix 秒)
	LastSign int64 `json:"last_sign"`
	// 提前多少秒提醒
	RemindTimeSeconds int64 `json:"remind_time_seconds"`
	// 上次发送提醒时对应的签到周期，用于避免重复提醒
	LastRemind int64 `json:"last_remind"`
	// 是否已经触发过，触发后进入静默状态直到用户重新签到
	Triggered bool `json:"triggered"`
	// 上次触发时间
	LastTrigger int64 `json:"last_trigger"`
	CreateAt    int64 `json:"create_at"`
}

// NextSignTime 下次必须完成签到的时间点
func (t *Timer) NextSignTime() int64 {
	return t.LastSign + t.SignDerationSeconds
}

// RemindAt 应该发送提醒的时间点
func (t *Timer) RemindAt() int64 {
	return t.NextSignTime() - t.RemindTimeSeconds
}

// TimerRequest 创建/编辑 Timer 的请求体
type TimerRequest struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	Enabled             *bool  `json:"enabled"`
	SignDerationSeconds int64  `json:"sign_deration_seconds"`
	RemindTimeSeconds   int64  `json:"remind_time_seconds"`
}

// Validate 校验 Timer 请求参数
func (r *TimerRequest) Validate() error {
	if r.Name == "" {
		return ErrValidation("Timer 名称不能为空")
	}
	if len(r.Name) > 64 {
		return ErrValidation("Timer 名称过长")
	}
	if r.SignDerationSeconds < 60 {
		return ErrValidation("签到周期不能小于 60 秒")
	}
	if r.RemindTimeSeconds < 0 {
		return ErrValidation("提醒提前量不能为负数")
	}
	if r.RemindTimeSeconds >= r.SignDerationSeconds {
		return ErrValidation("提醒提前量必须小于签到周期")
	}
	return nil
}
