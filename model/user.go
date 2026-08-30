package model

import "fmt"

type UserRegistryReuest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// NotifyConfigRequest 修改提醒渠道配置的请求
type NotifyConfigRequest struct {
	DingTalkConfig *string `json:"dingtalk_config"`
}

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	CreateAt int64  `json:"create_at"`
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"-"`
	APIKey   string `gorm:"index;size:67" json:"api_key"`
	// 管理员可以修改系统配置。第一个注册的用户自动成为管理员。
	IsAdmin bool `gorm:"default:false" json:"is_admin"`
	// 钉钉机器人配置(JSON 字符串)，用于接收提醒
	DingTalkConfig *string `json:"dingtalk_config"`
}

func (u *User) IsValid() error {
	if u.Username == "" || u.Password == "" {
		return fmt.Errorf("用户名或密码不能为空")
	}
	if len(u.Username) > 32 || len(u.Username) < 2 {
		return fmt.Errorf("用户名过长或过短")
	}
	if len(u.Password) > 64 || len(u.Password) < 6 {
		return fmt.Errorf("密码过长或过短")
	}
	return nil
}
