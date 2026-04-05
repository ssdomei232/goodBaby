package onebot

import (
	"encoding/json"
	"strconv"

	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
	napcat "github.com/ssdomei232/napcat-http-go-sdk"
)

// 从 Rule 中获取 OneBotConfig
func getOneBotConfigFromRule(rule *model.Rule) *OneBotConfig {
	var oneBotConfig OneBotConfig
	err := json.Unmarshal([]byte(rule.ConfigJson), &oneBotConfig)
	if err != nil {
		return nil
	}
	return &oneBotConfig
}

// 从 Rule 中获取 OneBotAccount
func getOneBotAccountFromRule(rule *model.Rule) (*OneBotAccount, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	var oneBotAccountEntry model.Account
	gormDB.Where("id = ?", rule.AccountID).First(&oneBotAccountEntry)

	var oneBotAccount OneBotAccount
	err = json.Unmarshal([]byte(oneBotAccountEntry.Config), &oneBotAccount)
	if err != nil {
		return nil, err
	}

	return &oneBotAccount, nil
}

// 发送一条消息到OneBot
func sendOneBotMsg(account *OneBotAccount, groupID int, msg string) error {
	client := napcat.NewClient(account.Token, account.URL)

	err := client.SendGroupMsg(strconv.Itoa(groupID), msg)
	return err
}
