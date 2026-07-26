package bilibili

import (
	"encoding/json"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

func newClient(rawCookies string) *bilibili.Client {
	client := bilibili.New()
	client.SetRawCookies(rawCookies)
	return client
}

// 获取 Bilibili Client
func getBiliClient(rule *model.Rule) (*bilibili.Client, error) {
	var accountConfig BiliAccount
	if err := db.LoadAccountConfig(rule.AccountID, &accountConfig); err != nil {
		return nil, err
	}

	return newClient(accountConfig.RawCookies), nil
}

// 获取 Bilibili Dynamic Config
func getBiliDynamicConfig(rule *model.Rule) (*BiliDynamicConfig, error) {
	var biliDynamicConfig BiliDynamicConfig
	if err := json.Unmarshal([]byte(rule.ConfigJson), &biliDynamicConfig); err != nil {
		return nil, err
	}

	return &biliDynamicConfig, nil
}
