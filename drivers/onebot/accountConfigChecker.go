package onebot

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/ssdomei232/goodBaby/internal/meta"
)

// OneBotAccountConfigValidator OneBot账号配置验证器
type OneBotAccountConfigValidator struct{}

func (v *OneBotAccountConfigValidator) GetType() string {
	return AccountType
}

func (v *OneBotAccountConfigValidator) Validate(config string) error {
	cfg, err := parseAccount(config)
	if err != nil {
		return err
	}

	if cfg.Token == "" {
		return fmt.Errorf("OneBot账号配置中Token不能为空")
	}
	if cfg.URL == "" {
		return fmt.Errorf("OneBot账号配置中URL不能为空")
	}

	parsed, err := url.Parse(cfg.URL)
	if err != nil || parsed.Host == "" || !strings.HasPrefix(parsed.Scheme, "http") {
		return fmt.Errorf("OneBot账号配置中URL不合法，应形如 http://localhost:3000")
	}

	return nil
}

// Test 调用 get_login_info 验证地址与 Token
func (v *OneBotAccountConfigValidator) Test(config string) error {
	cfg, err := parseAccount(config)
	if err != nil {
		return err
	}

	var resp apiResponse
	if err := newClient(cfg).DoRequest("POST", "/get_login_info", map[string]any{}, &resp); err != nil {
		return fmt.Errorf("连接 OneBot 失败: %w", err)
	}
	return resp.err()
}

func (v *OneBotAccountConfigValidator) Meta() meta.AccountMeta {
	return meta.AccountMeta{
		Type:        AccountType,
		Label:       "OneBot (QQ)",
		Description: "对接 NapCat / go-cqhttp 等 OneBot HTTP 服务，用于发送 QQ 消息。",
		Fields: []meta.Field{
			{Key: "url", Label: "HTTP 服务地址", Type: meta.FieldString, Required: true, Placeholder: "http://localhost:3000"},
			{Key: "token", Label: "Access Token", Type: meta.FieldPassword, Required: true, Secret: true},
		},
	}
}

func parseAccount(config string) (*OneBotAccount, error) {
	var cfg OneBotAccount
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return nil, fmt.Errorf("解析OneBot账号配置失败: %v", err)
	}
	return &cfg, nil
}
