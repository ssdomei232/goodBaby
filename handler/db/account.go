package db

import (
	"fmt"

	"github.com/ssdomei232/goodBaby/model"
)

// GetAccount 读取一个账号，返回数据库里的原始记录。
//
// 规则执行需要的账号凭据由上层(handler/runner)在这里取好再传给 driver，
// driver 自己不碰数据库。accountID 为 0 表示该规则不需要账号，返回 (nil, nil)。
func GetAccount(accountID uint) (*model.Account, error) {
	if accountID == 0 {
		return nil, nil
	}

	gormDB, err := GetGormDB()
	if err != nil {
		return nil, err
	}

	var account model.Account
	if err := gormDB.First(&account, accountID).Error; err != nil {
		return nil, fmt.Errorf("获取账号(ID: %d)失败: %w", accountID, err)
	}
	return &account, nil
}
