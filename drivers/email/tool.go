package email

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wneessen/go-mail"
)

// ParseAccountConfig 解析邮箱账号配置
func ParseAccountConfig(config string) (*EmailAccountConfig, error) {
	if config == "" {
		return nil, fmt.Errorf("解析邮箱账号配置失败: 规则没有关联账号")
	}

	var account EmailAccountConfig
	if err := json.Unmarshal([]byte(config), &account); err != nil {
		return nil, fmt.Errorf("解析邮箱账号配置失败: %v", err)
	}
	return &account, nil
}

// ParseRuleConfig 解析邮件规则配置
func ParseRuleConfig(configJSON string) (*EmailRule, error) {
	var config EmailRule
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析邮件规则配置失败: %v", err)
	}
	return &config, nil
}

// newSMTPClient 按账号配置里的加密方式创建 SMTP 客户端
func newSMTPClient(cfg *EmailAccountConfig) (*mail.Client, error) {
	options := []mail.Option{
		mail.WithPort(cfg.Port),
		mail.WithUsername(cfg.Username),
		mail.WithPassword(cfg.Password),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
	}

	switch cfg.SecurityOrDefault() {
	case SecuritySSL:
		options = append(options, mail.WithSSL())
	case SecuritySTARTTLS:
		options = append(options, mail.WithTLSPolicy(mail.TLSMandatory))
	case SecurityNone:
		options = append(options, mail.WithTLSPolicy(mail.NoTLS))
	}

	client, err := mail.NewClient(cfg.SMTPServer, options...)
	if err != nil {
		return nil, fmt.Errorf("创建邮件客户端失败: %w", err)
	}
	return client, nil
}

// buildMessage 组装一封邮件
func buildMessage(cfg *EmailAccountConfig, address, title, body string) (*mail.Msg, error) {
	message := mail.NewMsg()
	if err := message.From(cfg.FromOrDefault()); err != nil {
		return nil, fmt.Errorf("设置发件人失败: %w", err)
	}
	if err := message.To(address); err != nil {
		return nil, fmt.Errorf("设置收件人失败: %w", err)
	}
	message.Subject(title)
	message.SetBodyString(mail.TypeTextPlain, body)
	return message, nil
}

// sendMail 给一个地址发送一封邮件
func sendMail(ctx context.Context, cfg *EmailAccountConfig, config *EmailRule, address string) error {
	client, err := newSMTPClient(cfg)
	if err != nil {
		return err
	}

	message, err := buildMessage(cfg, address, config.Title, config.Msg)
	if err != nil {
		return err
	}

	if err := client.DialAndSendWithContext(ctx, message); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}
	return nil
}
