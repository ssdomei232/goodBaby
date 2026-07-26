// Package scheduler 管理定时检查任务，支持运行时调整间隔
package scheduler

import (
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/ssdomei232/goodBaby/handler/checker"
)

var (
	mu      sync.Mutex
	runner  *cron.Cron
	entryID cron.EntryID
)

// Start 启动定时检查，interval 单位为分钟
func Start(intervalMinutes int) error {
	mu.Lock()
	defer mu.Unlock()

	if runner != nil {
		return fmt.Errorf("定时任务已经启动")
	}

	runner = cron.New()
	id, err := runner.AddFunc(spec(intervalMinutes), checker.CheckTimers)
	if err != nil {
		runner = nil
		return fmt.Errorf("注册定时任务失败: %w", err)
	}

	entryID = id
	runner.Start()
	return nil
}

// Reschedule 用新的间隔重排定时任务，管理员改配置后调用
func Reschedule(intervalMinutes int) error {
	mu.Lock()
	defer mu.Unlock()

	if runner == nil {
		return fmt.Errorf("定时任务尚未启动")
	}

	// cron 不支持原地改周期，只能撤掉旧的再加一个
	runner.Remove(entryID)
	id, err := runner.AddFunc(spec(intervalMinutes), checker.CheckTimers)
	if err != nil {
		return fmt.Errorf("注册定时任务失败: %w", err)
	}

	entryID = id
	return nil
}

// Stop 停止定时任务，等待正在执行的任务结束
func Stop() {
	mu.Lock()
	defer mu.Unlock()

	if runner == nil {
		return
	}
	<-runner.Stop().Done()
	runner = nil
}

func spec(intervalMinutes int) string {
	if intervalMinutes <= 0 {
		intervalMinutes = 10
	}
	return fmt.Sprintf("@every %dm", intervalMinutes)
}
