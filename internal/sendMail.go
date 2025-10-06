package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/ssdomei232/goodBaby/configs"
	"github.com/wneessen/go-mail"
)

const (
	maxRetries = 10
	retryDelay = 10 * time.Second
)

// 向邮件列表发送死亡通告
func SendMail() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
		return
	}

	stopMsg := fmt.Sprintf("%s\n\n%s\n\n本消息由自动程序发送(摇篮系统)", GetBasicInfo(), config.MailContent)

	for _, address := range config.MailList {
		sendMailMsgWithRetry(address, stopMsg, config.MailTitle)
	}
}

func sendMailMsgWithRetry(address string, msg string, title string) {
	var lastErr error

	// 重试循环
	for i := 0; i <= maxRetries; i++ {
		if i > 0 {
			log.Printf("第 %d 次重试发送邮件给 %s", i, address)
		}

		err := sendMailMsg(address, msg, title)
		if err == nil {
			// 发送成功
			fmt.Printf("邮件成功发送给 %s\n", address)
			return
		}

		lastErr = err
		log.Printf("发送邮件给 %s 失败: %v", address, err)

		// 如果不是最后一次重试，则等待后重试
		if i < maxRetries {
			time.Sleep(retryDelay * time.Duration(i+1)) // 递增延迟
		}
	}

	// 所有重试都失败
	log.Printf("发送邮件给 %s 经过 %d 次重试后仍然失败: %v", address, maxRetries, lastErr)
}

func sendMailMsg(address string, msg string, title string) error {
	config, err := configs.GetConfig()
	if err != nil {
		return fmt.Errorf("获取配置文件失败: %v", err)
	}

	client, err := mail.NewClient(
		config.SMTPConfig.Host,
		mail.WithPort(config.SMTPConfig.Port),
		mail.WithSSL(),
		mail.WithUsername(config.SMTPConfig.User),
		mail.WithPassword(config.SMTPConfig.Password),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	)
	if err != nil {
		return fmt.Errorf("创建邮件客户端失败: %v", err)
	}

	// 创建邮件
	message := mail.NewMsg()
	if err := message.From(config.SMTPConfig.User); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	if err := message.To(address); err != nil {
		return fmt.Errorf("设置收件人失败: %v", err)
	}

	message.Subject(title)
	message.SetBodyString(mail.TypeTextPlain, msg)

	// 发送邮件
	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}
