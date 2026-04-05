package account

import (
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// 获取账号相关的规则列表
func getRulesByAccountID(accountID uint, uid uint) ([]*model.Rule, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	var rules []*model.Rule
	result := gormDB.Where("account_id = ? AND uid = ?", accountID, uid).Find(&rules)
	if result.Error != nil {
		return nil, result.Error
	}
	return rules, nil
}
