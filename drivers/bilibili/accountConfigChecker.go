package bilibili

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/internal/meta"
)

// BilibiliAccountConfigValidator B站账号配置验证器
type BilibiliAccountConfigValidator struct{}

func (v *BilibiliAccountConfigValidator) GetType() string {
	return AccountType
}

func (v *BilibiliAccountConfigValidator) Validate(config string) error {
	cfg, err := parseAccount(config)
	if err != nil {
		return err
	}

	if cfg.RawCookies == "" {
		return fmt.Errorf("B站账号配置中 cookies 不能为空")
	}

	return nil
}

// Test 用 cookie 拉一次账号信息，验证 cookie 是否还有效
func (v *BilibiliAccountConfigValidator) Test(config string) error {
	cfg, err := parseAccount(config)
	if err != nil {
		return err
	}

	client := newClient(cfg.RawCookies)
	info, err := client.GetAccountInformation()
	if err != nil {
		return fmt.Errorf("B站 cookie 校验失败: %w", err)
	}
	if info == nil || info.Mid == 0 {
		return fmt.Errorf("B站 cookie 无效或已过期")
	}
	return nil
}

func (v *BilibiliAccountConfigValidator) Meta() meta.AccountMeta {
	return meta.AccountMeta{
		Type:        AccountType,
		Label:       "哔哩哔哩",
		Description: "使用浏览器中的完整 Cookie 登录 B 站，用于发送动态。",
		Docs:        "docs/bilibili-config.md",
		Fields: []meta.Field{
			{
				Key:         "raw_cookies",
				Label:       "Cookies",
				Type:        meta.FieldTextarea,
				Required:    true,
				Secret:      true,
				Placeholder: "SESSDATA=xxx; bili_jct=xxx; DedeUserID=xxx",
				Help:        "登录 bilibili.com 后从浏览器开发者工具复制完整的 Cookie 字符串",
			},
		},
	}
}

func parseAccount(config string) (*BiliAccount, error) {
	var cfg BiliAccount
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return nil, fmt.Errorf("解析B站账号配置失败: %v", err)
	}
	return &cfg, nil
}
