package bilibili

import (
	"encoding/json"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// 获取 Bilibili Client
func GetBiliClient(rule *model.Rule) (*bilibili.Client, error) {
	client := bilibili.New()
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	// 1. 通过 account_id 获取到对应的 account 配置
	var biliAccount model.Account
	result := gormDB.Where("id = ?", rule.AccountID).First(&biliAccount)
	if result.Error != nil {
		return nil, result.Error
	}

	// 2. 通过 account 配置中的 config 字段获取到 CookiesString
	var biliAccountConfig BiliAccount
	err = json.Unmarshal([]byte(biliAccount.Config), &biliAccountConfig)
	if err != nil {
		return nil, err
	}

	client.SetRawCookies(biliAccountConfig.RawCookies)

	return client, nil
}
