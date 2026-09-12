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

// NewTask 组装执行一条规则所需的输入：规则本体 + 从数据库取出的关联账号。
//
// 数据库访问只发生在这一层，driver 拿到的 task 已经是自包含的。
func NewTask(rule *model.Rule) (*model.RuleTask, error) {
	account, err := db.GetAccount(rule.AccountID)
	if err != nil {
		return nil, err
	}
	return &model.RuleTask{Rule: rule, Account: account}, nil
}

// ExecuteRule 执行单条规则并记录执行日志
func ExecuteRule(rule *model.Rule, trigger string) error {
	ctx, cancel := retry.ExecutionContext()
	defer cancel()
	return ExecuteRuleWithContext(ctx, rule, trigger)
}

// ExecuteRuleWithContext 使用给定 context 执行规则，手动测试时可传入较短的超时
func ExecuteRuleWithContext(ctx context.Context, rule *model.Rule, trigger string) error {
	task, err := NewTask(rule)
	return execTask(ctx, task, err, &model.ExecutionLog{
		UID:      rule.UID,
		RuleID:   rule.ID,
		RuleName: rule.Name,
		RuleType: rule.Type,
		TimerID:  rule.TimerID,
		Trigger:  trigger,
	})
}

// ExecuteGatewayRuleWithContext 执行一条消息网关规则，日志里带上网关信息。
//
// 与定时器规则走同一条执行链路，区别只有触发来源与日志字段。
func ExecuteGatewayRuleWithContext(ctx context.Context, rule *model.GatewayRule, trigger string) error {
	task, err := NewTask(rule.AsRule())
	return execTask(ctx, task, err, &model.ExecutionLog{
		UID:           rule.UID,
		RuleName:      rule.Name,
		RuleType:      rule.Type,
		GatewayID:     rule.GatewayID,
		GatewayRuleID: rule.ID,
		Trigger:       trigger,
	})
}

// execTask 执行规则并写入执行日志。
//
// taskErr 是组装输入时的错误(例如账号已被删除)：这时不会走到 driver，
// 但同样要留下一条失败日志，方便用户在 WebUI 上看到原因。
func execTask(ctx context.Context, task *model.RuleTask, taskErr error, entry *model.ExecutionLog) error {
	if taskErr != nil {
		return record(entry, taskErr)
	}

	log.Printf("执行规则 [%s] %s (触发来源: %s)", entry.RuleType, entry.RuleName, entry.Trigger)
	return record(entry, GetGlobalExecutorRegistry().Execute(ctx, task))
}

// record 写入执行日志并返回原始错误
func record(entry *model.ExecutionLog, err error) error {
	entry.Success = err == nil
	entry.Message = "执行成功"
	if err != nil {
		entry.Message = err.Error()
		log.Printf("执行规则失败 [%s] %s: %v", entry.RuleType, entry.RuleName, err)
	}

	logstore.Record(entry)
	return err
}
