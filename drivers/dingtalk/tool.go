package dingtalk

import (
	"encoding/json"

	"github.com/ssdomei232/goodBaby/model"
)

func getDingTalkConfigFromUser(user *model.User) *DingTalkConfig {
	if user.DingTalkConfig == nil {
		return nil
	}

	var config DingTalkConfig
	err := json.Unmarshal([]byte(*user.DingTalkConfig), &config)
	if err != nil {
		return nil
	}
	return &config
}
