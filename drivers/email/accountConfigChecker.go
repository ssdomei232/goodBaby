package email

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"

	"github.com/ssdomei232/goodBaby/internal/meta"
	"github.com/ssdomei232/goodBaby/internal/retry"
)

// EmailAccountConfigValidator 邮箱账号配置验证器
type EmailAccountConfigValidator struct{}

func (v *EmailAccountConfigValidator) GetType() string {
	return AccountType
}

func (v *EmailAccountConfigValidator) Validate(config string) error {
	cfg, err := parseAccountConfig(config)
	if err != nil {
		return err
	}

	if cfg.SMTPServer == "" {
		return fmt.Errorf("邮箱账号配置中 SMTP 服务器不能为空")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("邮箱账号配置中 SMTP 端口不合法")
	}
	if cfg.Username == "" {
		return fmt.Errorf("邮箱账号配置中用户名不能为空")
	}
	if cfg.Password == "" {
		return fmt.Errorf("邮箱账号配置中密码不能为空")
	}
	switch cfg.SecurityOrDefault() {
	case SecuritySSL, SecuritySTARTTLS, SecurityNone:
	default:
		return fmt.Errorf("不支持的加密方式: %s", cfg.Security)
	}
	if cfg.From != "" {
		if _, err := mail.ParseAddress(cfg.From); err != nil {
			return fmt.Errorf("发件人地址不合法: %s", cfg.From)
		}
	}
	if cfg.TestDestination != "" {
		if _, err := mail.ParseAddress(cfg.TestDestination); err != nil {
			return fmt.Errorf("测试收件地址不合法: %s", cfg.TestDestination)
		}
	}

	return nil
}

// Test 连接 SMTP 服务器并完成认证；填写了测试收件地址时会真的发一封测试邮件
func (v *EmailAccountConfigValidator) Test(config string) error {
	cfg, err := parseAccountConfig(config)
	if err != nil {
		return err
	}

	client, err := newSMTPClient(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), retry.TestTimeout)
	defer cancel()

	if cfg.TestDestination == "" {
		if err := client.DialWithContext(ctx); err != nil {
			return fmt.Errorf("连接 SMTP 服务器失败: %w", err)
		}
		return client.Close()
	}

	message, err := buildMessage(cfg, cfg.TestDestination, "goodBaby 测试邮件", "这是一封来自 goodBaby 的测试邮件，收到即表示邮箱账号配置正确。")
	if err != nil {
		return err
	}
	if err := client.DialAndSendWithContext(ctx, message); err != nil {
		return fmt.Errorf("发送测试邮件失败: %w", err)
	}
	return nil
}

func (v *EmailAccountConfigValidator) Meta() meta.AccountMeta {
	return meta.AccountMeta{
		Type:        AccountType,
		Label:       "邮箱 (SMTP)",
		Description: "配置一个用于发信的 SMTP 账号。",
		Docs:        "docs/email-config.md",
		Fields: []meta.Field{
			{Key: "smtp_server", Label: "SMTP 服务器", Type: meta.FieldString, Required: true, Placeholder: "smtp.example.com"},
			{Key: "port", Label: "端口", Type: meta.FieldNumber, Required: true, Default: 465, Help: "SSL 一般为 465，STARTTLS 一般为 587"},
			{Key: "security", Label: "加密方式", Type: meta.FieldString, Required: false, Default: SecuritySSL, Help: "可选 ssl / starttls / none，留空按 ssl 处理"},
			{Key: "username", Label: "用户名", Type: meta.FieldString, Required: true, Placeholder: "you@example.com"},
			{Key: "password", Label: "密码 / 授权码", Type: meta.FieldPassword, Required: true, Secret: true},
			{Key: "from", Label: "发件人地址", Type: meta.FieldString, Help: "留空则使用用户名作为发件人"},
			{Key: "test_destination", Label: "测试收件地址", Type: meta.FieldString, Help: "填写后点击“测试”会真的发送一封测试邮件"},
		},
	}
}

func parseAccountConfig(config string) (*EmailAccountConfig, error) {
	var cfg EmailAccountConfig
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return nil, fmt.Errorf("解析邮箱账号配置失败: %v", err)
	}
	return &cfg, nil
}
