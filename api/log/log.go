// Package log 提供规则执行日志的查询接口
package log

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/response"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

const (
	defaultPageSize = 20
	maxPageSize     = 200
)

// HandleGetLogs 分页查询当前用户的执行日志
func HandleGetLogs(c *gin.Context) {
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

	page := intQuery(c, "page", 1, 1, 1<<20)
	pageSize := intQuery(c, "page_size", defaultPageSize, 1, maxPageSize)

	query := gormDB.Model(&model.ExecutionLog{}).Where("uid = ?", userInfo.ID)
	if ruleID := c.Query("rule_id"); ruleID != "" {
		query = query.Where("rule_id = ?", ruleID)
	}
	if success := c.Query("success"); success == "true" || success == "false" {
		query = query.Where("success = ?", success == "true")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.ServerError(c, "获取执行日志失败")
		return
	}

	logs := []model.ExecutionLog{}
	if err := query.Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		response.ServerError(c, "获取执行日志失败")
		return
	}

	response.OK(c, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"items":     logs,
	})
}

// HandleClearLogs 清空当前用户的执行日志
func HandleClearLogs(c *gin.Context) {
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

	if err := gormDB.Where("uid = ?", userInfo.ID).Delete(&model.ExecutionLog{}).Error; err != nil {
		response.ServerError(c, "清空执行日志失败")
		return
	}

	response.OK(c, "已清空执行日志")
}

func intQuery(c *gin.Context, key string, fallback, min, max int) int {
	raw := c.Query(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < min {
		return fallback
	}
	if value > max {
		return max
	}
	return value
}
