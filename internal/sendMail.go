package internal

import (
	"fmt"
	"log"

	"github.com/ssdomei232/goodBaby/configs"
	"github.com/wneessen/go-mail"
)

func SendMail() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
	}
	for _, address := range config.MailList {
		sendMailMsg(address, config.MailContent)
	}
}

func sendMailMsg(address string, msg string) {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
	}

	mailMsg := fmt.Sprintf("Hello,\nThis is a test email sent from Go!%s", msg)

	client, err := mail.NewClient(
		config.SMTPConfig.Host,
		mail.WithPort(config.SMTPConfig.Port),
		mail.WithSSL(),
		mail.WithUsername(config.SMTPConfig.User),
		mail.WithPassword(config.SMTPConfig.Password),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	)
	if err != nil {
		log.Println(err)
		return
	}

	// 创建邮件
	message := mail.NewMsg()
	message.From(config.SMTPConfig.User)
	message.To(address)
	message.Subject(config.MailTitle)
	message.SetBodyString(mail.TypeTextPlain, mailMsg)

	// 发送邮件
	if err := client.DialAndSend(message); err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Email sent successfully!")
}
