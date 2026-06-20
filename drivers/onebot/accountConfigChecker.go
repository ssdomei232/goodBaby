package onebot

import (
	"encoding/json"
	"fmt"
)

// AccountConfigValidator OneBot账号配置验证器接口
type OneBotAccountConfigValidator struct{}

func (v *OneBotAccountConfigValidator) GetType() string {
	return "onebot"
}

func (v *OneBotAccountConfigValidator) Validate(config string) error {
	var cfg OneBotAccount
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return fmt.Errorf("解析OneBot账号配置失败: %v", err)
	}

	if cfg.Token == "" {
		return fmt.Errorf("OneBot账号配置中Token不能为空")
	}
	if cfg.URL == "" {
		return fmt.Errorf("OneBot账号配置中URL不能为空")
	}

	return nil
}
