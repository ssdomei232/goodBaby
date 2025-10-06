package configs

import (
	"encoding/json"
	"os"
)

type Config struct {
	SignalSecret       string       `json:"signal_secret"`
	Debug              bool         `json:"debug"`
	DisconnectDuration int          `json:"disconnect_duration"` // Hours
	CatBotUrl          string       `json:"cat_bot_url"`
	CatBotKey          string       `json:"cat_bot_key"`
	QQSendGroup        []int        `json:"qq_send_group"`
	QQMsg              string       `json:"qq_msg"`
	MailList           []string     `json:"mail_list"`
	MailTitle          string       `json:"mail_title"`
	SMTPConfig         SMTPConfig   `json:"smtp_config"`
	MailContent        string       `json:"mail_content"`
	GithubConfig       GithubConfig `json:"github_config"`
	BiliMsg            string       `json:"bili_msg"`
}

type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"pass"`
}

type GithubConfig struct {
	Owner string   `json:"owner"`
	Repos []string `json:"repos"`
	Token string   `json:"token"`
}

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
