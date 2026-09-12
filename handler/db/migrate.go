package db

import (
	"log"

	"github.com/ssdomei232/goodBaby/model"
	"gorm.io/gorm"
)

// legacyRule 是迁移前的 rules 表结构，只用来读出已经废弃的 gateway_id
type legacyRule struct {
	ID         uint
	UID        uint
	GatewayID  uint
	Name       string
	AccountID  uint
	Type       string
	ConfigJson string
	Enabled    bool
	CreateAt   int64
}

// migrateGatewayRules 把历史数据里混在 rules 表中的消息网关规则搬到 gateway_rules 表。
//
// v2.1 之前消息网关规则和定时器规则共用 rules 表，靠 gateway_id 区分；
// 现在两者分表存储，这里做一次性搬迁，避免老用户升级后规则凭空消失。
func migrateGatewayRules(gormDB *gorm.DB) error {
	if !gormDB.Migrator().HasColumn("rules", "gateway_id") {
		return nil
	}

	var legacy []legacyRule
	if err := gormDB.Table("rules").Where("gateway_id <> 0").Find(&legacy).Error; err != nil {
		return err
	}
	if len(legacy) == 0 {
		return nil
	}

	return gormDB.Transaction(func(tx *gorm.DB) error {
		for _, rule := range legacy {
			moved := model.GatewayRule{
				UID:        rule.UID,
				GatewayID:  rule.GatewayID,
				Name:       rule.Name,
				AccountID:  rule.AccountID,
				Type:       rule.Type,
				ConfigJson: rule.ConfigJson,
				Enabled:    rule.Enabled,
				CreateAt:   rule.CreateAt,
			}
			if err := tx.Create(&moved).Error; err != nil {
				return err
			}
			if err := tx.Exec("DELETE FROM rules WHERE id = ?", rule.ID).Error; err != nil {
				return err
			}
		}
		log.Printf("已迁移 %d 条消息网关规则到 gateway_rules 表", len(legacy))
		return nil
	})
}
