package checker

import (
	"log"
	"time"

	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/handler/reminder"
	"github.com/ssdomei232/goodBaby/handler/runner"
	"github.com/ssdomei232/goodBaby/model"
)

// 定时检查所有 timers 是否需要被触发
//
// 同时向即将被触发的 timers 发送通知
func CheckTimers() {
	gormDB, err := db.GetGormDB()
	if err != nil {
		log.Printf("failed to get gorm db: %v", err)
		return
	}

	// 1. 获取所有 timers
	var timers []model.Timer
	if err := gormDB.Find(&timers).Error; err != nil {
		log.Printf("failed to get timers: %v", err)
		return
	}

	// 2. 检查每个 timer 是否需要被触发
	for _, timer := range timers {
		// 计算下次触发时间
		nextSignTime := timer.LastSign + timer.SignDerationSeconds

		// 计算提醒时间
		remindTime := nextSignTime - timer.RemindTimeSeconds

		// 获取当前时间
		currentTime := time.Now().Unix()

		if currentTime >= remindTime && currentTime < nextSignTime {
			// 发送通知
			reminder.Reminder(&timer)
		}

		if currentTime >= nextSignTime {
			// 触发 timer
			runner.Runner(&timer)
		}
	}
}
