package rule

import (
	"fmt"
	"strconv"

	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/internal/meta"
	"github.com/ssdomei232/goodBaby/internal/ruleConfigChecker"
	"github.com/ssdomei232/goodBaby/model"
)

func parseID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func boolOr(v *bool, fallback bool) bool {
	if v == nil {
		return fallback
	}
	return *v
}

// findRule 查找属于该用户的规则
func findRule(ruleID, uid uint) (*model.Rule, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	var target model.Rule
	if err := gormDB.Where("id = ? AND uid = ?", ruleID, uid).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

// DeleteRuleByID 根据 ID 删除规则
func DeleteRuleByID(id uint, uid uint) error {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return err
	}

	return gormDB.Where("id = ? AND uid = ?", id, uid).Delete(&model.Rule{}).Error
}

// checkRuleConfigAccountAndTimerExist 检查关联的 Timer 与账号是否存在且归属当前用户
//
// requiredAccountType 为空表示该规则类型不需要账号。
func checkRuleConfigAccountAndTimerExist(rule *model.Rule, requiredAccountType string) error {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return err
	}

	var count int64
	if err := gormDB.Model(&model.Timer{}).
		Where("id = ? AND uid = ?", rule.TimerID, rule.UID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return model.ErrValidation("关联的 Timer 不存在")
	}

	if requiredAccountType == "" {
		// 该类型不需要账号，忽略前端可能传来的 account_id
		rule.AccountID = 0
		return nil
	}

	if rule.AccountID == 0 {
		return model.ErrValidation("该规则类型必须关联一个账号")
	}

	var account model.Account
	if err := gormDB.Where("id = ? AND uid = ?", rule.AccountID, rule.UID).First(&account).Error; err != nil {
		return model.ErrValidation("关联的账号不存在")
	}
	if account.Type != requiredAccountType {
		return model.ErrValidation(fmt.Sprintf("规则需要 %s 类型的账号，但关联的是 %s", requiredAccountType, account.Type))
	}

	return nil
}

// maskRule 掩码规则配置中的敏感字段
func maskRule(target model.Rule) model.Rule {
	if ruleMeta, ok := ruleConfigChecker.InitValidatorRegistry().MetaOf(target.Type); ok {
		target.ConfigJson = meta.Mask(target.ConfigJson, ruleMeta.Fields)
	}
	return target
}

func maskRules(rules []model.Rule) []model.Rule {
	masked := make([]model.Rule, 0, len(rules))
	for _, r := range rules {
		masked = append(masked, maskRule(r))
	}
	return masked
}

// getRuleOwnerUID 获取规则所属用户的 UID
func getRuleOwnerUID(ruleID uint) (uint, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return 0, err
	}

	var rule model.Rule
	if err := gormDB.Select("uid").Where("id = ?", ruleID).First(&rule).Error; err != nil {
		return 0, err
	}

	return rule.UID, nil
}
