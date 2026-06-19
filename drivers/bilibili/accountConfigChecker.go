package bilibili

import (
	"encoding/json"
	"fmt"
)

// BilibiliAccountConfigValidator B站账号配置验证器
type BilibiliAccountConfigValidator struct{}

func (v *BilibiliAccountConfigValidator) GetType() string {
	return "bilibili-account"
}

func (v *BilibiliAccountConfigValidator) Validate(config string) error {
	var cfg BiliAccount
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return fmt.Errorf("解析B站账号配置失败: %v", err)
	}

	if cfg.RawCookies == "" {
		return fmt.Errorf("B站账号配置中 cookies 不能为空")
	}

	return nil
}
