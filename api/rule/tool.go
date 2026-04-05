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
