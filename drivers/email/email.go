package email

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/ssdomei232/goodBaby/model"
)

// SendMail 向规则配置里的所有地址发送邮件。
//
// 每个地址单独重试，互不阻塞；全部失败/部分失败都会汇总成错误返回，
// 由 runner 记录到执行日志中。
func SendMail(ctx context.Context, rule *model.Rule) error {
	emailRule, err := GetEmailRuleFromRule(rule)
	if err != nil {
		return fmt.Errorf("获取邮件规则配置失败: %w", err)
	}

	accountConfig, err := GetEmailAccountFromRule(rule)
	if err != nil {
		return fmt.Errorf("获取邮件账户配置失败: %w", err)
	}

	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		fails []string
	)

	for _, destination := range emailRule.Destinations {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			if err := sendMailMsgWithRetry(ctx, accountConfig, emailRule, address); err != nil {
				mu.Lock()
				fails = append(fails, fmt.Sprintf("%s: %v", address, err))
				mu.Unlock()
			}
		}(destination)
	}
	wg.Wait()

	if len(fails) > 0 {
		return fmt.Errorf("%d/%d 封邮件发送失败: %s",
			len(fails), len(emailRule.Destinations), strings.Join(fails, "; "))
	}
	return nil
}
