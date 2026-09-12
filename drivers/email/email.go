package email

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/ssdomei232/goodBaby/internal/retry"
)

// SendMail 用给定 SMTP 账号向规则配置里的所有地址发送邮件。
//
// 每个地址单独重试、互不阻塞；全部失败/部分失败都会汇总成错误返回，
// 由 runner 记录到执行日志中。
func SendMail(ctx context.Context, account *EmailAccountConfig, config *EmailRule) error {
	var (
		wg    sync.WaitGroup
		mu    sync.Mutex
		fails []string
	)

	for _, destination := range config.Destinations {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			if err := retry.Do(ctx, func() error {
				return sendMail(ctx, account, config, address)
			}); err != nil {
				mu.Lock()
				fails = append(fails, fmt.Sprintf("%s: %v", address, err))
				mu.Unlock()
			}
		}(destination)
	}
	wg.Wait()

	if len(fails) > 0 {
		// 并发发送的顺序不确定，排序后错误信息才是稳定的
		sort.Strings(fails)
		return fmt.Errorf("%d/%d 封邮件发送失败: %s",
			len(fails), len(config.Destinations), strings.Join(fails, "; "))
	}
	return nil
}
