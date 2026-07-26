package timer

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/response"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
	"gorm.io/gorm"
)

// HandleGetAllTimers 获取当前用户的所有 Timer
func HandleGetAllTimers(c *gin.Context) {
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
	if err := gormDB.Where("uid = ?", userInfo.ID).Order("id DESC").Find(&timers).Error; err != nil {
		response.ServerError(c, "获取 Timer 失败")
		return
	}

	response.OK(c, timers)
}

// HandleGetTimer 获取单个 Timer 及其关联规则数量
func HandleGetTimer(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	timerID, err := parseID(c.Param("timerID"))
	if err != nil {
		response.BadRequest(c, "Timer ID 格式错误")
		return
	}

	timer, err := findTimer(timerID, userInfo.ID)
	if err != nil {
		response.NotFound(c, "Timer 不存在")
		return
	}

	ruleCount, err := countRules(timerID, userInfo.ID)
	if err != nil {
		response.ServerError(c, "获取关联规则失败")
		return
	}

	response.OK(c, gin.H{"timer": timer, "rule_count": ruleCount})
}

// HandleCreateTimer 创建 Timer
func HandleCreateTimer(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	var req model.TimerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入参数错误")
		return
	}

	if err := req.Validate(); err != nil {
		response.FromError(c, err, "创建 Timer 失败")
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	now := time.Now().Unix()
	timer := model.Timer{
		UID:                 userInfo.ID,
		Name:                req.Name,
		Description:         req.Description,
		Enabled:             boolOr(req.Enabled, true),
		SignDerationSeconds: req.SignDerationSeconds,
		RemindTimeSeconds:   req.RemindTimeSeconds,
		// 创建即视为完成一次签到，从当前时间开始计时
		LastSign: now,
		CreateAt: now,
	}

	if err := gormDB.Create(&timer).Error; err != nil {
		response.ServerError(c, "创建 Timer 失败")
		return
	}

	response.OK(c, timer)
}

// HandleEditTimer 编辑 Timer
func HandleEditTimer(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	timerID, err := parseID(c.Param("timerID"))
	if err != nil {
		response.BadRequest(c, "Timer ID 格式错误")
		return
	}

	timer, err := findTimer(timerID, userInfo.ID)
	if err != nil {
		response.NotFound(c, "Timer 不存在")
		return
	}

	var req model.TimerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入参数错误")
		return
	}

	if err := req.Validate(); err != nil {
		response.FromError(c, err, "更新 Timer 失败")
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	timer.Name = req.Name
	timer.Description = req.Description
	timer.Enabled = boolOr(req.Enabled, timer.Enabled)
	timer.SignDerationSeconds = req.SignDerationSeconds
	timer.RemindTimeSeconds = req.RemindTimeSeconds

	if err := gormDB.Save(timer).Error; err != nil {
		response.ServerError(c, "更新 Timer 失败")
		return
	}

	response.OK(c, timer)
}

// HandleSignTimer 签到：重置计时并解除已触发状态
func HandleSignTimer(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	timerID, err := parseID(c.Param("timerID"))
	if err != nil {
		response.BadRequest(c, "Timer ID 格式错误")
		return
	}

	timer, err := findTimer(timerID, userInfo.ID)
	if err != nil {
		response.NotFound(c, "Timer 不存在")
		return
	}

	if err := signTimer(timer); err != nil {
		response.ServerError(c, "签到失败")
		return
	}

	response.OK(c, timer)
}

// HandleSignAll 一键签到当前用户的所有启用中的 Timer
func HandleSignAll(c *gin.Context) {
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

	now := time.Now().Unix()
	result := gormDB.Model(&model.Timer{}).
		Where("uid = ? AND enabled = ?", userInfo.ID, true).
		Updates(map[string]any{
			"last_sign":   now,
			"last_remind": 0,
			"triggered":   false,
		})
	if result.Error != nil {
		response.ServerError(c, "签到失败")
		return
	}

	response.OK(c, gin.H{"signed": result.RowsAffected, "last_sign": now})
}

// HandleDeleteTimer 删除 Timer
//
// Timer 被删除后，挂在它下面的规则就没有触发来源了，因此一并删除。
func HandleDeleteTimer(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	timerID, err := parseID(c.Param("timerID"))
	if err != nil {
		response.BadRequest(c, "Timer ID 格式错误")
		return
	}

	if _, err := findTimer(timerID, userInfo.ID); err != nil {
		response.NotFound(c, "Timer 不存在")
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	err = gormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("timer_id = ? AND uid = ?", timerID, userInfo.ID).Delete(&model.Rule{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND uid = ?", timerID, userInfo.ID).Delete(&model.Timer{}).Error
	})
	if err != nil {
		response.ServerError(c, "删除 Timer 失败")
		return
	}

	response.OK(c, "Timer 删除成功")
}

// HandleCheckDeleteTimer 删除前查询会被一并删除的规则
func HandleCheckDeleteTimer(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	timerID, err := parseID(c.Param("timerID"))
	if err != nil {
		response.BadRequest(c, "Timer ID 格式错误")
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	rules := []model.Rule{}
	if err := gormDB.Where("timer_id = ? AND uid = ?", timerID, userInfo.ID).Find(&rules).Error; err != nil {
		response.ServerError(c, "获取相关规则失败")
		return
	}

	response.OK(c, rules)
}
