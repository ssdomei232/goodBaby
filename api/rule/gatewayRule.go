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
	"github.com/ssdomei232/goodBaby/internal/gateway"
	"github.com/ssdomei232/goodBaby/internal/meta"
	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

// HandleGetAllGatewayRules 获取用户的消息网关规则，支持按 gateway_id 过滤
//
// 消息网关规则和定时器规则是两张表，互不影响，因此这里单独一组接口。
func HandleGetAllGatewayRules(c *gin.Context) {
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
	if raw := c.Query("gateway_id"); raw != "" {
		gatewayID, err := parseID(raw)
		if err != nil {
			response.BadRequest(c, "gateway_id 格式错误")
			return
		}
		query = query.Where("gateway_id = ?", gatewayID)
	}

	rules := []model.GatewayRule{}
	if err := query.Order("id DESC").Find(&rules).Error; err != nil {
		response.ServerError(c, "获取网关规则失败")
		return
	}

	response.OK(c, maskGatewayRules(rules))
}

// HandleCreateGatewayRule 创建消息网关规则
func HandleCreateGatewayRule(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	var req model.GatewayRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入参数错误")
		return
	}

	newRule := model.GatewayRule{
		UID:        userInfo.ID,
		Name:       req.Name,
		GatewayID:  req.GatewayID,
		AccountID:  req.AccountID,
		Type:       req.Type,
		ConfigJson: req.ConfigJson,
		Enabled:    boolOr(req.Enabled, true),
		CreateAt:   time.Now().Unix(),
	}

	if err := validateGatewayRule(&req, &newRule); err != nil {
		response.FromError(c, err, "创建网关规则失败")
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	if err := gormDB.Create(&newRule).Error; err != nil {
		response.ServerError(c, "创建网关规则失败")
		return
	}

	response.OK(c, maskGatewayRule(newRule))
}

// HandleEditGatewayRule 编辑消息网关规则
func HandleEditGatewayRule(c *gin.Context) {
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

	existing, err := findGatewayRule(ruleID, userInfo.ID)
	if err != nil {
		response.NotFound(c, "网关规则不存在")
		return
	}

	var req model.GatewayRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入参数错误")
		return
	}

	updated := *existing
	updated.Name = req.Name
	updated.Type = req.Type
	updated.GatewayID = req.GatewayID
	updated.AccountID = req.AccountID
	updated.Enabled = boolOr(req.Enabled, existing.Enabled)
	// 前端提交的敏感字段可能是掩码占位符，用旧配置补回
	updated.ConfigJson = unmaskRuleConfig(req.Type, req.ConfigJson, existing.ConfigJson)

	req.ConfigJson = updated.ConfigJson
	if err := validateGatewayRule(&req, &updated); err != nil {
		response.FromError(c, err, "更新网关规则失败")
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	if err := gormDB.Save(&updated).Error; err != nil {
		response.ServerError(c, "更新网关规则失败")
		return
	}

	response.OK(c, maskGatewayRule(updated))
}

// HandleDeleteGatewayRule 删除消息网关规则
func HandleDeleteGatewayRule(c *gin.Context) {
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

	if _, err := findGatewayRule(ruleID, userInfo.ID); err != nil {
		response.NotFound(c, "网关规则不存在")
		return
	}

	if err := DeleteGatewayRuleByID(ruleID, userInfo.ID); err != nil {
		response.ServerError(c, "删除网关规则失败")
		return
	}

	response.OK(c, "网关规则删除成功")
}

// HandleTestGatewayRule 立即执行一次网关规则，用于验证配置
//
// 与投递不同：这里用规则里保存的配置，不覆盖外部消息。
func HandleTestGatewayRule(c *gin.Context) {
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

	target, err := findGatewayRule(ruleID, userInfo.ID)
	if err != nil {
		response.NotFound(c, "网关规则不存在")
		return
	}

	// 用一条示例消息走一次真实的投递覆盖，让用户看到规则最终会发出什么内容
	if configJSON, applied, err := gateway.ApplyMessage(target.ConfigJson, gateway.TestMessage); err == nil && applied {
		target.ConfigJson = configJSON
	}

	ctx, cancel := context.WithTimeout(context.Background(), retry.TestTimeout)
	defer cancel()

	if err := runner.ExecuteGatewayRuleWithContext(ctx, target, model.TriggerManual); err != nil {
		response.BadRequest(c, fmt.Sprintf("规则执行失败: %s", err.Error()))
		return
	}

	response.OK(c, "规则执行成功")
}

// validateGatewayRule 校验网关规则的通用字段、关联对象与类型专属配置
func validateGatewayRule(req *model.GatewayRuleRequest, rule *model.GatewayRule) error {
	if err := req.Validate(); err != nil {
		return err
	}

	if err := checkGatewayExists(rule.GatewayID, rule.UID); err != nil {
		return err
	}

	ruleMeta, err := ruleMetaOf(req.Type)
	if err != nil {
		return err
	}

	// 标题/内容由投递请求提供，页面不填，这里补占位值后再校验
	rule.ConfigJson = meta.FillGatewayMessages(rule.ConfigJson, ruleMeta.Fields)
	if err := validateRuleConfig(req.Type, rule.ConfigJson); err != nil {
		return err
	}

	accountID, err := checkAccount(rule.AccountID, rule.UID, ruleMeta.AccountType)
	if err != nil {
		return err
	}
	rule.AccountID = accountID

	return nil
}
