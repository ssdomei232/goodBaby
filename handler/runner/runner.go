package runner

import (
	"log"

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
		go executeRule(&rule)
	}
}

// executeRule 执行单个规则
func executeRule(rule *model.Rule) {
	executorRegistry := GetGlobalExecutorRegistry()

	if err := executorRegistry.Execute(rule); err != nil {
		log.Printf("执行规则失败 [ID: %d, Type: %s]: %v", rule.ID, rule.Type, err)
	}
}
