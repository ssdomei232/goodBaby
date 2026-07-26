package checker

import (
	"log"
	"time"

	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/handler/reminder"
	"github.com/ssdomei232/goodBaby/handler/runner"
	"github.com/ssdomei232/goodBaby/model"
)

// CheckTimers 定时检查所有 timers 是否需要被触发
//
// 同时向即将到期的 timers 对应的用户发送提醒
func CheckTimers() {
	gormDB, err := db.GetGormDB()
	if err != nil {
		log.Printf("failed to get gorm db: %v", err)
		return
	}

	// 1. 获取所有启用中的 timers
	var timers []model.Timer
	if err := gormDB.Where("enabled = ?", true).Find(&timers).Error; err != nil {
		log.Printf("failed to get timers: %v", err)
		return
	}

	currentTime := time.Now().Unix()

	// 2. 检查每个 timer 是否需要被提醒或触发
	for i := range timers {
		timer := timers[i]

		// 已触发过的 timer 在用户重新签到之前不再处理
		if timer.Triggered {
			continue
		}

		nextSignTime := timer.NextSignTime()
		remindTime := timer.RemindAt()

		if currentTime >= nextSignTime {
			runner.Runner(&timer)
			continue
		}

		// 每个签到周期只提醒一次
		if currentTime >= remindTime && timer.LastRemind < remindTime {
			reminder.Reminder(&timer)
			if err := gormDB.Model(&model.Timer{}).Where("id = ?", timer.ID).
				Update("last_remind", currentTime).Error; err != nil {
				log.Printf("更新 Timer(ID: %d) 提醒时间失败: %v", timer.ID, err)
			}
		}
	}
}
