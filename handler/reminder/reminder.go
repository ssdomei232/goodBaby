package reminder

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ssdomei232/goodBaby/drivers/dingtalk"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/handler/logstore"
	"github.com/ssdomei232/goodBaby/model"
)

// Reminder 在 checker 检查到 timer 即将到期时，向用户发送提醒
func Reminder(timer *model.Timer) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		log.Printf("failed to get gorm db: %v", err)
		return
	}

	// 1. 获取 timer 对应的 user
	var user model.User
	if err := gormDB.Where("id = ?", timer.UID).First(&user).Error; err != nil {
		log.Printf("failed to get user: %v", err)
		return
	}

	// 2. 向用户发送提醒
	deadline := time.Unix(timer.NextSignTime(), 0).Format("2006-01-02 15:04:05")
	title := fmt.Sprintf("goodBaby: %s 即将到期", timer.Name)
	msg := fmt.Sprintf("### %s\n\n定时器 **%s** 将在 %s 到期。\n\n请及时签到，否则关联的规则会被执行。", title, timer.Name, deadline)

	err = dingtalk.SendDingTalkMsg(&user, title, msg)
	if errors.Is(err, dingtalk.ErrNoUserConfig) {
		// 用户没有配置提醒渠道，不算失败
		return
	}

	message := "提醒发送成功"
	if err != nil {
		message = err.Error()
		log.Printf("发送提醒失败 [Timer ID: %d]: %v", timer.ID, err)
	}

	logstore.Record(&model.ExecutionLog{
		UID:      timer.UID,
		TimerID:  timer.ID,
		RuleName: timer.Name,
		RuleType: "reminder",
		Trigger:  model.TriggerRemind,
		Success:  err == nil,
		Message:  message,
	})
}
