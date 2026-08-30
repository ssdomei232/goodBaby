package account

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/response"
	"github.com/ssdomei232/goodBaby/api/rule"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/internal/accountConfigChecker"
	"github.com/ssdomei232/goodBaby/internal/meta"
	"github.com/ssdomei232/goodBaby/model"
)

// HandleAddAccount 添加账号
func HandleAddAccount(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	var req model.AccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入参数错误")
		return
	}

	if err := req.Validate(); err != nil {
		response.FromError(c, err, "添加账号失败")
		return
	}

	// 之前这里完全没有走校验器，任何 JSON 都能存进去，直到规则触发才暴露问题
	registry := accountConfigChecker.InitValidatorRegistry()
	if err := registry.Validate(req.Type, req.Config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	newAccount := model.Account{
		UID:      userInfo.ID,
		Name:     req.Name,
		Type:     req.Type,
		Config:   req.Config,
		CreateAt: time.Now().Unix(),
	}

	if err := gormDB.Create(&newAccount).Error; err != nil {
		response.ServerError(c, "添加账号失败")
		return
	}

	response.OK(c, maskAccount(newAccount))
}

// HandleGetAllAccounts 获取用户的所有账号
func HandleGetAllAccounts(c *gin.Context) {
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
	if accountType := c.Query("type"); accountType != "" {
		query = query.Where("type = ?", accountType)
	}

	accounts := []model.Account{}
	if err := query.Order("id DESC").Find(&accounts).Error; err != nil {
		response.ServerError(c, "获取账号失败")
		return
	}

	response.OK(c, maskAccounts(accounts))
}

// HandleEditAccount 编辑账号
func HandleEditAccount(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	accountID, err := parseID(c.Param("accountID"))
	if err != nil {
		response.BadRequest(c, "账号ID格式错误")
		return
	}

	existing, err := findAccount(accountID, userInfo.ID)
	if err != nil {
		response.NotFound(c, "账号不存在")
		return
	}

	var req model.AccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入参数错误")
		return
	}

	if err := req.Validate(); err != nil {
		response.FromError(c, err, "更新账号失败")
		return
	}

	if req.Type != existing.Type {
		response.BadRequest(c, "不支持修改账号类型，请新建账号")
		return
	}

	registry := accountConfigChecker.InitValidatorRegistry()
	// 前端提交的密码等字段可能是掩码占位符，用旧配置补回
	config := meta.Unmask(req.Config, existing.Config, registry.FieldsOf(req.Type))

	if err := registry.Validate(req.Type, config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	existing.Name = req.Name
	existing.Config = config
	if err := gormDB.Save(existing).Error; err != nil {
		response.ServerError(c, "更新账号失败")
		return
	}

	response.OK(c, maskAccount(*existing))
}

// HandleTestAccount 测试账号凭据是否可用
func HandleTestAccount(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	accountID, err := parseID(c.Param("accountID"))
	if err != nil {
		response.BadRequest(c, "账号ID格式错误")
		return
	}

	existing, err := findAccount(accountID, userInfo.ID)
	if err != nil {
		response.NotFound(c, "账号不存在")
		return
	}

	registry := accountConfigChecker.InitValidatorRegistry()
	if err := registry.Test(existing.Type, existing.Config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "账号可用")
}

// HandleCheckDeleteAccount 检查删除账号请求
//
// 删除账号会同时删除相关规则，先请求该接口获取受影响的规则
func HandleCheckDeleteAccount(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	accountID, err := parseID(c.Param("accountID"))
	if err != nil {
		response.BadRequest(c, "账号ID格式错误")
		return
	}

	rules, err := getRulesByAccountID(accountID, userInfo.ID)
	if err != nil {
		response.ServerError(c, "获取相关规则失败")
		return
	}

	response.OK(c, rules)
}

// HandleDeleteAccount 删除账号及其关联规则
func HandleDeleteAccount(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	accountID, err := parseID(c.Param("accountID"))
	if err != nil {
		response.BadRequest(c, "账号ID格式错误")
		return
	}

	// 检查账号是否归属请求用户
	ownerUID, err := getAccountOwnerUID(accountID)
	if err != nil {
		response.ServerError(c, "获取账号所属用户失败")
		return
	}

	if ownerUID != userInfo.ID {
		response.Forbidden(c, "无权限操作该账号")
		return
	}

	// 删除相关规则
	rules, err := getRulesByAccountID(accountID, userInfo.ID)
	if err != nil {
		response.ServerError(c, "获取相关规则失败")
		return
	}
	for _, oneRule := range rules {
		if err := rule.DeleteRuleByID(oneRule.ID, userInfo.ID); err != nil {
			response.ServerError(c, "删除相关规则失败")
			return
		}
	}

	// 删除账号
	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}
	if err := gormDB.Where("id = ? AND uid = ?", accountID, userInfo.ID).
		Delete(&model.Account{}).Error; err != nil {
		response.ServerError(c, "删除账号失败")
		return
	}

	response.OK(c, "删除账号成功")
}
