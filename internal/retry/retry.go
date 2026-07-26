// Package retry 封装规则执行使用的指数退避重试
package retry

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/ssdomei232/goodBaby/configs"
)

// 手动测试时使用的超时，避免用户在 WebUI 上等待数小时
const TestTimeout = 30 * time.Second

// Do 在 ctx 允许的时间内按指数退避执行 op
//
// 超时后返回最后一次真实的业务错误，而不是干巴巴的 context deadline exceeded。
func Do(ctx context.Context, op func() error) error {
	var lastErr error
	_, err := backoff.Retry(ctx, func() (struct{}, error) {
		e := op()
		if e != nil {
			lastErr = e
		}
		return struct{}{}, e
	}, backoff.WithBackOff(backoff.NewExponentialBackOff()))

	if err != nil && lastErr != nil {
		return lastErr
	}
	return err
}

// ExecutionTimeout 返回配置中的规则执行超时时间
func ExecutionTimeout() time.Duration {
	config, err := configs.GetConfig()
	if err != nil || config.TimeoutDurationHours <= 0 {
		return 6 * time.Hour
	}
	return time.Duration(config.TimeoutDurationHours) * time.Hour
}

// ExecutionContext 返回一个带规则执行超时的 context
func ExecutionContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), ExecutionTimeout())
}
