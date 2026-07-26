package account

import (
	"strconv"

	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/internal/accountConfigChecker"
	"github.com/ssdomei232/goodBaby/internal/meta"
	"github.com/ssdomei232/goodBaby/model"
)

func parseID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// findAccount 查找属于该用户的账号
func findAccount(accountID, uid uint) (*model.Account, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	var account model.Account
	if err := gormDB.Where("id = ? AND uid = ?", accountID, uid).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

// getRulesByAccountID 获取账号相关的规则列表
func getRulesByAccountID(accountID uint, uid uint) ([]*model.Rule, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	rules := []*model.Rule{}
	if err := gormDB.Where("account_id = ? AND uid = ?", accountID, uid).Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// maskAccount 掩码账号配置中的敏感字段，避免 Cookie/密码原样返回给前端
func maskAccount(account model.Account) model.Account {
	fields := accountConfigChecker.InitValidatorRegistry().FieldsOf(account.Type)
	account.Config = meta.Mask(account.Config, fields)
	return account
}

func maskAccounts(accounts []model.Account) []model.Account {
	masked := make([]model.Account, 0, len(accounts))
	for _, a := range accounts {
		masked = append(masked, maskAccount(a))
	}
	return masked
}
