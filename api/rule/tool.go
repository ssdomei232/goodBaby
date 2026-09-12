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

// findRule 查找属于该用户的定时器规则
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

// findGatewayRule 查找属于该用户的消息网关规则
func findGatewayRule(ruleID, uid uint) (*model.GatewayRule, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	var target model.GatewayRule
	if err := gormDB.Where("id = ? AND uid = ?", ruleID, uid).First(&target).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

// DeleteRuleByID 根据 ID 删除定时器规则
func DeleteRuleByID(id uint, uid uint) error {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return err
	}

	return gormDB.Where("id = ? AND uid = ?", id, uid).Delete(&model.Rule{}).Error
}

// DeleteGatewayRuleByID 根据 ID 删除消息网关规则
func DeleteGatewayRuleByID(id uint, uid uint) error {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return err
	}

	return gormDB.Where("id = ? AND uid = ?", id, uid).Delete(&model.GatewayRule{}).Error
}

// checkTimerExists 检查关联的 Timer 是否存在且归属当前用户
func checkTimerExists(timerID, uid uint) error {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return err
	}

	var count int64
	if err := gormDB.Model(&model.Timer{}).Where("id = ? AND uid = ?", timerID, uid).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return model.ErrValidation("关联的 Timer 不存在")
	}
	return nil
}

// checkGatewayExists 检查关联的消息网关是否存在且归属当前用户
func checkGatewayExists(gatewayID, uid uint) error {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return err
	}

	var count int64
	if err := gormDB.Model(&model.MessageGateway{}).Where("id = ? AND uid = ?", gatewayID, uid).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return model.ErrValidation("关联的消息网关不存在")
	}
	return nil
}

// checkAccount 校验规则关联的账号存在、归属当前用户且类型匹配。
//
// requiredType 为空表示该规则类型不需要账号，此时返回 0，调用方据此清掉 account_id。
func checkAccount(accountID, uid uint, requiredType string) (uint, error) {
	if requiredType == "" {
		return 0, nil
	}
	if accountID == 0 {
		return 0, model.ErrValidation("该规则类型必须关联一个账号")
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		return 0, err
	}

	var account model.Account
	if err := gormDB.Where("id = ? AND uid = ?", accountID, uid).First(&account).Error; err != nil {
		return 0, model.ErrValidation("关联的账号不存在")
	}
	if account.Type != requiredType {
		return 0, model.ErrValidation(
			fmt.Sprintf("规则需要 %s 类型的账号，但关联的是 %s", requiredType, account.Type))
	}
	return accountID, nil
}

// ruleMetaOf 取规则类型的元数据，不支持的类型返回校验错误
func ruleMetaOf(ruleType string) (meta.RuleMeta, error) {
	ruleMeta, ok := ruleConfigChecker.InitValidatorRegistry().MetaOf(ruleType)
	if !ok {
		return meta.RuleMeta{}, model.ErrValidation(fmt.Sprintf("不支持的规则类型: %s", ruleType))
	}
	return ruleMeta, nil
}

// validateRuleConfig 校验类型专属的规则配置
func validateRuleConfig(ruleType, configJSON string) error {
	if err := ruleConfigChecker.InitValidatorRegistry().Validate(ruleType, configJSON); err != nil {
		return model.ErrValidation(fmt.Sprintf("规则配置验证失败: %s", err.Error()))
	}
	return nil
}

// maskRule 掩码定时器规则配置中的敏感字段
func maskRule(target model.Rule) model.Rule {
	target.ConfigJson = maskConfig(target.Type, target.ConfigJson)
	return target
}

func maskRules(rules []model.Rule) []model.Rule {
	masked := make([]model.Rule, 0, len(rules))
	for _, r := range rules {
		masked = append(masked, maskRule(r))
	}
	return masked
}

// maskGatewayRule 掩码消息网关规则配置中的敏感字段
func maskGatewayRule(target model.GatewayRule) model.GatewayRule {
	target.ConfigJson = maskConfig(target.Type, target.ConfigJson)
	return target
}

func maskGatewayRules(rules []model.GatewayRule) []model.GatewayRule {
	masked := make([]model.GatewayRule, 0, len(rules))
	for _, r := range rules {
		masked = append(masked, maskGatewayRule(r))
	}
	return masked
}

// maskConfig 掩码规则配置中的敏感字段
func maskConfig(ruleType, configJSON string) string {
	if ruleMeta, ok := ruleConfigChecker.InitValidatorRegistry().MetaOf(ruleType); ok {
		return meta.Mask(configJSON, ruleMeta.Fields)
	}
	return configJSON
}

// getRuleOwnerUID 获取定时器规则所属用户的 UID
func getRuleOwnerUID(ruleID uint) (uint, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return 0, err
	}

	var target model.Rule
	if err := gormDB.Select("uid").Where("id = ?", ruleID).First(&target).Error; err != nil {
		return 0, err
	}
	return target.UID, nil
}

// getGatewayRuleOwnerUID 获取消息网关规则所属用户的 UID
func getGatewayRuleOwnerUID(ruleID uint) (uint, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return 0, err
	}

	var target model.GatewayRule
	if err := gormDB.Select("uid").Where("id = ?", ruleID).First(&target).Error; err != nil {
		return 0, err
	}
	return target.UID, nil
}
