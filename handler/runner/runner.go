package runner

import (
	"log"

	"github.com/ssdomei232/goodBaby/drivers/bilibili"
	"github.com/ssdomei232/goodBaby/drivers/email"
	"github.com/ssdomei232/goodBaby/drivers/github"
	"github.com/ssdomei232/goodBaby/drivers/onebot"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// runner 执行需要被执行的 Rule
func Runner(timer *model.Timer) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		log.Printf("failed to get gorm db: %v", err)
		return
	}

	// 1. 获取所有需要执行的 Rule
	var rules []model.Rule
	if err := gormDB.Where("timer_id = ?", timer.ID).Find(&rules).Error; err != nil {
		log.Printf("failed to get rules: %v", err)
		return
	}

	// 2. 执行每个 Rule
	for _, rule := range rules {
		// 执行 rule
		go executeRule(&rule)
	}
}

// super switch
func executeRule(rule *model.Rule) {
	switch rule.Type {
	case "bilibili-dynamic":
		bilibili.SendBiliDynamicMsg(rule)
	case "email":
		email.SendMail(rule)
	case "github-make-repo-public":
		github.MakeRepositoryPublic(rule)
	case "onebot":
		onebot.SendOneBotMsg(rule)
	}
}
