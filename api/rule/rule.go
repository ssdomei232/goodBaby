package rule

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/response"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/handler/runner"
	"github.com/ssdomei232/goodBaby/internal/meta"
	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/internal/ruleConfigChecker"
	"github.com/ssdomei232/goodBaby/model"
)

// HandleGetAllRules 获取用户的所有规则，支持按 timer_id 过滤
func HandleGetAllRules(c *gin.Context) {
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

	query := gormDB.Where("uid = ?", userInfo.ID)
	if raw := c.Query("timer_id"); raw != "" {
		timerID, err := parseID(raw)
		if err != nil {
			response.BadRequest(c, "timer_id 格式错误")
			return
		}
		query = query.Where("timer_id = ?", timerID)
	}

	rules := []model.Rule{}
	if err := query.Order("id DESC").Find(&rules).Error; err != nil {
		response.ServerError(c, "获取规则失败")
		return
	}

	response.OK(c, maskRules(rules))
}

// HandleCreateRule 创建新规则
func HandleCreateRule(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	var req model.RuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入参数错误")
		return
	}

	newRule := model.Rule{
		UID:        userInfo.ID,
		Name:       req.Name,
		TimerID:    req.TimerID,
		AccountID:  req.AccountID,
		Type:       req.Type,
		ConfigJson: req.ConfigJson,
		Enabled:    boolOr(req.Enabled, true),
		CreateAt:   time.Now().Unix(),
	}

	if err := validateRule(&req, &newRule); err != nil {
		response.FromError(c, err, "创建规则失败")
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	if err := gormDB.Create(&newRule).Error; err != nil {
		response.ServerError(c, "创建规则失败")
		return
	}

	response.OK(c, maskRule(newRule))
}

// HandleEditRule 编辑规则
func HandleEditRule(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	ruleID, err := parseID(c.Param("ruleID"))
	if err != nil {
		response.BadRequest(c, "规则 ID 格式错误")
		return
	}

	existing, err := findRule(ruleID, userInfo.ID)
	if err != nil {
		response.NotFound(c, "规则不存在")
		return
	}

	var req model.RuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入参数错误")
		return
	}

	updated := *existing
	updated.Name = req.Name
	updated.Type = req.Type
	updated.TimerID = req.TimerID
	updated.AccountID = req.AccountID
	updated.Enabled = boolOr(req.Enabled, existing.Enabled)
	// 前端提交的敏感字段可能是掩码占位符，用旧配置补回
	updated.ConfigJson = unmaskRuleConfig(req.Type, req.ConfigJson, existing.ConfigJson)

	req.ConfigJson = updated.ConfigJson
	if err := validateRule(&req, &updated); err != nil {
		response.FromError(c, err, "更新规则失败")
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	if err := gormDB.Save(&updated).Error; err != nil {
		response.ServerError(c, "更新规则失败")
		return
	}

	response.OK(c, maskRule(updated))
}

// HandleDeleteRule 删除规则
func HandleDeleteRule(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	ruleID, err := parseID(c.Param("ruleID"))
	if err != nil {
		response.BadRequest(c, "规则 ID 格式错误")
		return
	}

	if err := DeleteRuleByID(ruleID, userInfo.ID); err != nil {
		response.ServerError(c, "删除规则失败")
		return
	}

	response.OK(c, "规则删除成功")
}

// HandleTestRule 立即执行一次规则用于验证配置
//
// 使用较短的超时，避免在 WebUI 上等待数小时的指数退避。
func HandleTestRule(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	ruleID, err := parseID(c.Param("ruleID"))
	if err != nil {
		response.BadRequest(c, "规则 ID 格式错误")
		return
	}

	target, err := findRule(ruleID, userInfo.ID)
	if err != nil {
		response.NotFound(c, "规则不存在")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), retry.TestTimeout)
	defer cancel()

	if err := runner.ExecuteRuleWithContext(ctx, target, model.TriggerManual); err != nil {
		response.BadRequest(c, fmt.Sprintf("规则执行失败: %s", err.Error()))
		return
	}

	response.OK(c, "规则执行成功")
}

// validateRule 校验规则的通用字段、关联对象与类型专属配置
func validateRule(req *model.RuleRequest, rule *model.Rule) error {
	if err := req.Validate(); err != nil {
		return err
	}

	ruleMeta, ok := ruleConfigChecker.InitValidatorRegistry().MetaOf(req.Type)
	if !ok {
		return model.ErrValidation(fmt.Sprintf("不支持的规则类型: %s", req.Type))
	}

	// 检查关联的 Timer 与账号是否存在且属于当前用户
	if err := checkRuleConfigAccountAndTimerExist(rule, ruleMeta.AccountType); err != nil {
		return err
	}

	if err := ruleConfigChecker.InitValidatorRegistry().Validate(req.Type, rule.ConfigJson); err != nil {
		return model.ErrValidation(fmt.Sprintf("规则配置验证失败: %s", err.Error()))
	}

	return nil
}

// unmaskRuleConfig 把提交上来的掩码字段还原成旧值
func unmaskRuleConfig(ruleType, newConfig, oldConfig string) string {
	ruleMeta, ok := ruleConfigChecker.InitValidatorRegistry().MetaOf(ruleType)
	if !ok {
		return newConfig
	}
	return meta.Unmask(newConfig, oldConfig, ruleMeta.Fields)
}
