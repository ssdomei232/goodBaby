package email

import (
	"encoding/json"
	"fmt"
)

// EmailAccountConfigValidator 邮箱账号配置验证器
type EmailAccountConfigValidator struct{}

func (v *EmailAccountConfigValidator) GetType() string {
	return "email-account"
}

func (v *EmailAccountConfigValidator) Validate(config string) error {
	// 邮箱账号配置验证逻辑
	var cfg EmailAccountConfig
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return fmt.Errorf("解析邮箱账号配置失败: %v", err)
	}

	if cfg.SMTPServer == "" {
		return fmt.Errorf("邮箱账号配置中 SMTP 服务器不能为空")
	}
	if cfg.Port == 0 {
		return fmt.Errorf("邮箱账号配置中 SMTP 端口不能为空")
	}
	if cfg.Username == "" {
		return fmt.Errorf("邮箱账号配置中用户名不能为空")
	}
	if cfg.Password == "" {
		return fmt.Errorf("邮箱账号配置中密码不能为空")
	}

	return nil
}
