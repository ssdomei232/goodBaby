package model

import (
	"fmt"

	"github.com/ssdomei232/goodBaby/handler/db"
)

type UserRegistryReuest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID             uint    `json:"id" gorm:"primaryKey"`
	CreateAt       int64   `json:"create_at"`
	Username       string  `json:"username"`
	Password       string  `json:"-"`
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

func (u *User) IsExist() bool {
	db, err := db.GetGormDB()
	if err != nil {
		return false
	}

	var count int64
	result := db.Where("username = ?", u.Username).Count(&count)
	if result.Error != nil {
		return false
	}

	return count > 0
}
