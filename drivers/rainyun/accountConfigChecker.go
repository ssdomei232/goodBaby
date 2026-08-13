package rainyun

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/internal/meta"
	rain "github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/common"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/user"
)

// RainyunAccountConfigValidator Rainyun账号配置验证器
type RainyunAccountConfigValidator struct{}

func (v *RainyunAccountConfigValidator) GetType() string {
	return AccountType
}

func (v *RainyunAccountConfigValidator) Validate(config string) error {
	// 解析配置
	cfg, err := parseAccount(config)
	if err != nil {
		return err
	}

	if cfg.APIKey == "" {
		return fmt.Errorf("雨云账号配置中APIKey不能为空")
	}

	return nil
}

func parseAccount(config string) (*RainyunAccount, error) {
	var cfg RainyunAccount
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return nil, fmt.Errorf("解析雨云账号配置失败: %v", err)
	}
	return &cfg, nil
}

func (v *RainyunAccountConfigValidator) Test(config string) error {
	cfg, err := parseAccount(config)
	if err != nil {
		return err
	}

	client := rain.NewClient(cfg.APIKey)
	userClient := user.Client{Client: client}
	_, err = userClient.GetUserInfo()

	if err != nil {
		return fmt.Errorf("测试雨云账号配置失败: %v", err)
	}

	return nil
}

func (v *RainyunAccountConfigValidator) Meta() meta.AccountMeta {
	return meta.AccountMeta{
		Type:        AccountType,
		Label:       "雨云",
		Description: "对接雨云，方便各位代理/运维删库跑路",
		Fields: []meta.Field{
			{Key: "api_key", Label: "API Key", Type: meta.FieldPassword, Required: true, Secret: true, Placeholder: "xxxxxxxxx"},
		},
	}
}
