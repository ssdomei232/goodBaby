package configs

import (
	"encoding/json"
	"os"
)

// type DingtalkBot struct {
// 	AccessToken string `json:"access_token"`
// 	Secret      string `json:"secret"`
// }

type Config struct {
	EnableRegistry       bool `json:"enable_registry"`
	TimeoutDurationHours int  `json:"timeout_duration_hours"`
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
