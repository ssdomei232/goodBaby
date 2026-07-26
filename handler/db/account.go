package db

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/model"
)

// LoadAccountConfig 读取账号并把 Config 字段反序列化到 out。
//
// 各个 driver 之前都各自实现了一遍这段逻辑，其中 github driver 甚至
// 直接把驱动内的结构体当成 gorm model 查询，这里统一收敛。
func LoadAccountConfig(accountID uint, out any) error {
	if accountID == 0 {
		return fmt.Errorf("该规则没有关联账号")
	}

	gormDB, err := GetGormDB()
	if err != nil {
		return err
	}

	var account model.Account
	if err := gormDB.First(&account, accountID).Error; err != nil {
		return fmt.Errorf("获取账号(ID: %d)失败: %w", accountID, err)
	}

	if err := json.Unmarshal([]byte(account.Config), out); err != nil {
		return fmt.Errorf("解析账号(ID: %d)配置失败: %w", accountID, err)
	}
	return nil
}
