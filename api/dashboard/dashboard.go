// Package dashboard 汇总首页需要的统计数据
package dashboard

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/response"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// HandleGetOverview 返回首页概览：各类对象数量、最紧急的 Timer、最近的执行日志
func HandleGetOverview(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	timers := []model.Timer{}
	if err := gormDB.Where("uid = ?", userInfo.ID).Find(&timers).Error; err != nil {
		response.ServerError(c, "获取统计数据失败")
		return
	}

	var ruleCount, accountCount int64
	gormDB.Model(&model.Rule{}).Where("uid = ?", userInfo.ID).Count(&ruleCount)
	gormDB.Model(&model.Account{}).Where("uid = ?", userInfo.ID).Count(&accountCount)

	recentLogs := []model.ExecutionLog{}
	gormDB.Where("uid = ?", userInfo.ID).Order("id DESC").Limit(10).Find(&recentLogs)

	now := time.Now().Unix()
	var (
		enabledCount   int
		triggeredCount int
		urgent         *model.Timer
	)

	for i := range timers {
		timer := timers[i]
		if !timer.Enabled {
			continue
		}
		enabledCount++
		if timer.Triggered {
			triggeredCount++
			continue
		}
		// 剩余时间最短的启用中 Timer
		if urgent == nil || timer.NextSignTime() < urgent.NextSignTime() {
			urgent = &timers[i]
		}
	}

	overview := gin.H{
		"timer_count":     len(timers),
		"enabled_timers":  enabledCount,
		"triggered_count": triggeredCount,
		"rule_count":      ruleCount,
		"account_count":   accountCount,
		"server_time":     now,
		"recent_logs":     recentLogs,
	}

	if urgent != nil {
		overview["urgent_timer"] = urgent
		overview["urgent_seconds_left"] = urgent.NextSignTime() - now
	}

	response.OK(c, overview)
}
