package reminder

import (
	"log"

	"github.com/ssdomei232/goodBaby/drivers/dingtalk"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// reminder 负责在 checker 检查到 timer 即将被触发时，向用户发送提醒
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
	dingtalk.SendDingTalkMsg(&user, "您的 goodbaby timer 即将被触发", "您的 goodbaby timer 即将被触发")
}
