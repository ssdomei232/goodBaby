package email

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
	"github.com/wneessen/go-mail"
)

// GetEmailAccountFromRule 从 Rule 中获取 EmailAccount 配置
func GetEmailAccountFromRule(rule *model.Rule) (*EmailAccountConfig, error) {
	var emailAccountConfig EmailAccountConfig
	if err := db.LoadAccountConfig(rule.AccountID, &emailAccountConfig); err != nil {
		return nil, err
	}
	return &emailAccountConfig, nil
}

// GetEmailRuleFromRule 从 Rule 中获取 EmailRule 配置
func GetEmailRuleFromRule(rule *model.Rule) (*EmailRule, error) {
	var emailRule EmailRule
	if err := json.Unmarshal([]byte(rule.ConfigJson), &emailRule); err != nil {
		return nil, err
	}
	return &emailRule, nil
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

func sendMailMsgWithRetry(ctx context.Context, cfg *EmailAccountConfig, rule *EmailRule, address string) error {
	return retry.Do(ctx, func() error {
		return sendMailMsg(ctx, cfg, rule, address)
	})
}

func sendMailMsg(ctx context.Context, cfg *EmailAccountConfig, rule *EmailRule, address string) error {
	client, err := newSMTPClient(cfg)
	if err != nil {
		return err
	}

	message, err := buildMessage(cfg, address, rule.Title, rule.Msg)
	if err != nil {
		return err
	}

	if err := client.DialAndSendWithContext(ctx, message); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}
	return nil
}
