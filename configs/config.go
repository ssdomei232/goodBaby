package configs

import (
	"encoding/json"
	"os"
)

type Config struct {
	SignalSecret       string       `json:"signal_secret"`
	Debug              bool         `json:"debug"`
	DisconnectDuration int          `json:"disconnect_duration"` // Hours
	BiliMsg            string       `json:"bili_msg"`
	Basic              Basic        `json:"basic"`
	MailConfig         MailConfig   `json:"mail_config"`
	GithubConfig       GithubConfig `json:"github_config"`
	DingtalkBot        DingtalkBot  `json:"dingtalk_bot"`
	OneBotConfig       OneBotConfig `json:"onebot_config"`
}

type MailConfig struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	User        string   `json:"user"`
	Password    string   `json:"pass"`
	MailList    []string `json:"mail_list"`
	MailTitle   string   `json:"mail_title"`
	MailContent string   `json:"mail_content"`
}

type GithubConfig struct {
	Owner string   `json:"owner"`
	Repos []string `json:"repos"`
	Token string   `json:"token"`
}

type Basic struct {
	Name      string `json:"name"`
	Age       int    `json:"age"`
	NickName  string `json:"nickname"`
	QQNumber  string `json:"qq_number"`
	CauseStop string `json:"cause_stop"`
}

type DingtalkBot struct {
	AccessToken string `json:"access_token"`
	Secret      string `json:"secret"`
}

type OneBotConfig struct {
	URL        string `json:"url"` // 例如 "http://localhost:5700"
	Token      string `json:"token"`
	SendGroups []int  `json:"send_groups"`
	Msg        string `json:"msg"`
}

// Get config from json file
func GetConfig() (config Config, err error) {
	content, err := os.ReadFile("config.json")
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
