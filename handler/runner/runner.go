package runner

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/handler/logstore"
	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

// Runner 执行某个 Timer 下所有启用的规则。
//
// 触发前先把 Timer 标记为已触发，避免规则执行期间(可能长达数小时)
// 被下一轮检查重复触发；用户重新签到后 Triggered 会被重置。
func Runner(timer *model.Timer) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		log.Printf("failed to get gorm db: %v", err)
		return
	}

	// 1. 获取所有需要执行的 Rule
	var rules []model.Rule
	if err := gormDB.Where("timer_id = ? AND enabled = ?", timer.ID, true).Find(&rules).Error; err != nil {
		log.Printf("failed to get rules: %v", err)
		return
	}

	// 2. 标记已触发
	now := time.Now().Unix()
	if err := gormDB.Model(&model.Timer{}).Where("id = ?", timer.ID).
		Updates(map[string]any{"triggered": true, "last_trigger": now}).Error; err != nil {
		log.Printf("标记 Timer(ID: %d) 已触发失败: %v", timer.ID, err)
		return
	}

	log.Printf("Timer %s (ID: %d) 已到期，开始执行 %d 条规则", timer.Name, timer.ID, len(rules))

	// 3. 并发执行每个 Rule
	var wg sync.WaitGroup
	for i := range rules {
		rule := rules[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			ExecuteRule(&rule, model.TriggerTimer)
		}()
	}
	wg.Wait()
}

// ExecuteRule 执行单个规则并记录执行日志
func ExecuteRule(rule *model.Rule, trigger string) error {
	ctx, cancel := retry.ExecutionContext()
	defer cancel()
	return ExecuteRuleWithContext(ctx, rule, trigger)
}

// ExecuteRuleWithContext 使用给定 context 执行规则，手动测试时可传入较短的超时
func ExecuteRuleWithContext(ctx context.Context, rule *model.Rule, trigger string) error {
	err := GetGlobalExecutorRegistry().Execute(ctx, rule)

	message := "执行成功"
	if err != nil {
		message = err.Error()
		log.Printf("执行规则失败 [ID: %d, Type: %s]: %v", rule.ID, rule.Type, err)
	}

	logstore.Record(&model.ExecutionLog{
		UID:      rule.UID,
		RuleID:   rule.ID,
		RuleName: rule.Name,
		RuleType: rule.Type,
		TimerID:  rule.TimerID,
		Trigger:  trigger,
		Success:  err == nil,
		Message:  message,
	})

	return err
}
