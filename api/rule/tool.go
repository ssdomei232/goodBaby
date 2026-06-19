package rule

import (
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// 根据 ID 删除规则
func DeleteRuleByID(id uint, uid uint) error {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return err
	}

	result := gormDB.Where("id = ? AND uid = ?", id, uid).Delete(&model.Rule{})
	return result.Error
}

// 检查关联账号和 Timer 是否存在(可以没有关联账号)
func checkRuleConfigAccountAndTimerExist(rule model.Rule) (exist bool, err error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return false, err
	}

	if rule.TimerID == 0 {
		return false, nil
	}

	var count int64
	if rule.AccountID != 0 {
		result := gormDB.Model(&model.Account{}).Where("id = ? AND uid = ?", rule.AccountID, rule.UID).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
		if count == 0 {
			return false, nil
		}
	}

	if rule.TimerID != 0 {
		result := gormDB.Model(&model.Timer{}).Where("id = ? AND uid = ?", rule.TimerID, rule.UID).Count(&count)
		if result.Error != nil {
			return false, result.Error
		}
		if count == 0 {
			return false, nil
		}
	}

	return true, nil
}
