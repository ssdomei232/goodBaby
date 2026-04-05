package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
	"github.com/wneessen/go-mail"
)

// 从 Rule 中获取 EmailAccount 配置
func GetEmailAccountFromRule(rule *model.Rule) (*EmailAccountConfig, error) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	// 1. 通过 account_id 获取到对应的 account 配置
	var emailAccount model.Account
	result := gormDB.Where("id = ?", rule.AccountID).First(&emailAccount)
	if result.Error != nil {
		return nil, result.Error
	}

	// 2. 通过 account 配置中的 config 字段获取到 EmailAccount
	var emailAccountConfig EmailAccountConfig
	err = json.Unmarshal([]byte(emailAccount.Config), &emailAccountConfig)
	if err != nil {
		return nil, err
	}

	return &emailAccountConfig, nil
}

// 从 Rule 中获取 EmailRule 配置
func GetEmailRuleFromRule(rule *model.Rule) (*EmailRule, error) {
	var emailRule EmailRule
	err := json.Unmarshal([]byte(rule.ConfigJson), &emailRule)
	if err != nil {
		return nil, err
	}
	return &emailRule, nil
}
func sendMailMsgWithRetry(rule *model.Rule, address string) {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置失败: %v", err)
		return
	}
	var timeout time.Duration = time.Duration(config.TimeoutDurationHours) * time.Hour

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	operation := func() (string, error) {
		err := sendMailMsg(address, rule)
		return "", err
	}

	_, err = backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		log.Printf("发送邮件失败: %v", err)
	}
}

func sendMailMsg(address string, rule *model.Rule) error {
	emailAccountConfig, err := GetEmailAccountFromRule(rule)
	if err != nil {
		log.Printf("获取邮件账户配置失败: %v", err)
		return fmt.Errorf("获取邮件账户配置失败: %v", err)
	}

	emailRule, err := GetEmailRuleFromRule(rule)
	if err != nil {
		log.Printf("获取邮件规则配置失败: %v", err)
		return fmt.Errorf("获取邮件规则配置失败: %v", err)
	}

	client, err := mail.NewClient(
		emailAccountConfig.SMTPServer,
		mail.WithPort(emailAccountConfig.Port),
		mail.WithSSL(),
		mail.WithUsername(emailAccountConfig.Username),
		mail.WithPassword(emailAccountConfig.Password),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	)
	if err != nil {
		return fmt.Errorf("创建邮件客户端失败: %v", err)
	}

	// 创建邮件
	message := mail.NewMsg()
	if err := message.From(emailAccountConfig.Username); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	if err := message.To(address); err != nil {
		return fmt.Errorf("设置收件人失败: %v", err)
	}

	message.Subject(emailRule.Title)
	message.SetBodyString(mail.TypeTextPlain, emailRule.Msg)

	// 发送邮件
	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}
