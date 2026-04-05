package email

import (
	"log"

	"github.com/ssdomei232/goodBaby/model"
)

// 向邮件列表发送消息
func SendMail(rule *model.Rule) {
	emailRule, err := GetEmailRuleFromRule(rule)
	if err != nil {
		log.Printf("获取邮件规则配置失败: %v", err)
		return
	}

	for _, destinations := range emailRule.Destinations {
		go sendMailMsgWithRetry(rule, destinations) // 使用 goroutine 发送邮件，避免指数退避阻塞其他地址
	}
}
