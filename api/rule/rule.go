package rule

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/internal/ruleConfigChecker"
	"github.com/ssdomei232/goodBaby/model"
)

// 获取用户的所有规则
func HandleGetAllRules(c *gin.Context) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "获取规则失败"})
		return
	}

	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		c.JSON(401, gin.H{"code": 401, "data": "获取用户信息失败"})
		return
	}

	var rules []model.Rule
	result := gormDB.Where("uid = ?", userInfo.ID).Find(&rules)
	if result.Error != nil {
		c.JSON(500, gin.H{"code": 500, "data": "获取规则失败"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": rules})
}

// 创建新规则
func HandleCreateRule(c *gin.Context) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "服务器内部错误"})
		return
	}

	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		c.JSON(401, gin.H{"code": 401, "data": "获取用户信息失败"})
		return
	}

	var newRule model.Rule
	if err := c.BindJSON(&newRule); err != nil {
		c.JSON(400, gin.H{"code": 400, "data": "输入参数错误"})
		return
	}

	if newRule.Type == "" {
		c.JSON(400, gin.H{"code": 400, "data": "规则类型不能为空"})
		return
	}

	if newRule.Name == "" {
		c.JSON(400, gin.H{"code": 400, "data": "规则名称不能为空"})
		return
	}

	// 检查关联账号和 Timer 是否存在(可以没有关联账号)
	if exist, err := checkRuleConfigAccountAndTimerExist(newRule); err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "服务器内部错误"})
		return
	} else if !exist {
		c.JSON(400, gin.H{"code": 400, "data": "关联的账户或Timer不存在"})
		return
	}

	// 规则校验
	validatorRegistry := ruleConfigChecker.InitValidatorRegistry()
	if err := validatorRegistry.Validate(newRule.Type, newRule.ConfigJson); err != nil {
		c.JSON(400, gin.H{"code": 400, "data": fmt.Sprintf("规则配置验证失败: %s", err.Error())})
		return
	}

	newRule.UID = userInfo.ID
	result := gormDB.Create(&newRule)
	if result.Error != nil {
		c.JSON(500, gin.H{"code": 500, "data": "创建规则失败"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": "规则创建成功"})
}

// 编辑规则
func HandleEditRule(c *gin.Context) {
	gormDB, err := db.GetGormDB()
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "服务器内部错误"})
		return
	}

	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		c.JSON(401, gin.H{"code": 401, "data": "获取用户信息失败"})
		return
	}

	ruleID, err := strconv.Atoi(c.Param("ruleID"))
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "data": "规则 ID 格式错误"})
		return
	}
	var rule model.Rule
	result := gormDB.Where("id = ? AND uid = ?", ruleID, userInfo.ID).First(&rule)
	if result.Error != nil {
		c.JSON(404, gin.H{"code": 404, "data": "规则不存在"})
		return
	}

	var updatedRule model.Rule
	if err := c.BindJSON(&updatedRule); err != nil {
		c.JSON(400, gin.H{"code": 400, "data": "输入参数错误"})
		return
	}

	if updatedRule.Type == "" {
		c.JSON(400, gin.H{"code": 400, "data": "规则类型不能为空"})
		return
	}

	if updatedRule.Name == "" {
		c.JSON(400, gin.H{"code": 400, "data": "规则名称不能为空"})
		return
	}

	// 检查关联账号和 Timer 是否存在(可以没有关联账号)
	if exist, err := checkRuleConfigAccountAndTimerExist(updatedRule); err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "服务器内部错误"})
		return
	} else if !exist {
		c.JSON(400, gin.H{"code": 400, "data": "关联的账户或Timer不存在"})
		return
	}

	// 规则校验
	validatorRegistry := ruleConfigChecker.InitValidatorRegistry()
	if err := validatorRegistry.Validate(updatedRule.Type, updatedRule.ConfigJson); err != nil {
		c.JSON(400, gin.H{"code": 400, "data": fmt.Sprintf("规则配置验证失败: %s", err.Error())})
		return
	}

	rule.Name = updatedRule.Name
	rule.Type = updatedRule.Type
	rule.ConfigJson = updatedRule.ConfigJson
	rule.AccountID = updatedRule.AccountID
	rule.TimerID = updatedRule.TimerID

	result = gormDB.Save(&rule)
	if result.Error != nil {
		c.JSON(500, gin.H{"code": 500, "data": "更新规则失败"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": "规则更新成功"})
}

// 根据删除规则
func HandleDeleteRule(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		c.JSON(401, gin.H{"code": 401, "data": "获取用户信息失败"})
		return
	}

	ruleID, err := strconv.Atoi(c.Param("ruleID"))
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "data": "规则 ID 格式错误"})
		return
	}

	err = DeleteRuleByID(uint(ruleID), userInfo.ID)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "删除规则失败"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": "规则删除成功"})
}
