package account

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/rule"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// 添加账号
func HandleAddAccount(c *gin.Context) {
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

	var newAccount model.Account
	if err := c.BindJSON(&newAccount); err != nil {
		c.JSON(400, gin.H{"code": 400, "data": "输入参数错误"})
		return
	}

	newAccount.UID = userInfo.ID
	result := gormDB.Create(&newAccount)
	if result.Error != nil {
		c.JSON(500, gin.H{"code": 500, "data": "添加账号失败"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": "添加账号成功"})
}

// 获取用户的所有账号
func HandleGetAllAccounts(c *gin.Context) {
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

	var accounts []model.Account
	result := gormDB.Where("uid = ?", userInfo.ID).Find(&accounts)
	if result.Error != nil {
		c.JSON(500, gin.H{"code": 500, "data": "获取账号失败"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": accounts})
}

// 检查删除账号请求
//
// 删除账号会同时删除相关规则，先请求该接口获取受影响的规则
func HandleCheckDeleteAccount(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		c.JSON(401, gin.H{"code": 401, "data": "获取用户信息失败"})
		return
	}
	accountID, err := strconv.Atoi(c.Param("accountID"))
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "data": "账号ID格式错误"})
		return
	}

	rules, err := getRulesByAccountID(uint(accountID), userInfo.ID)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "获取相关规则失败"})
		return
	}

	c.JSON(200, gin.H{"code": 200, "data": rules})
}

// 删除账号
func HandleDeleteAccount(c *gin.Context) {
	userInfo, err := user.GetUserInfoByGinCtx(c)
	if err != nil {
		c.JSON(401, gin.H{"code": 401, "data": "获取用户信息失败"})
		return
	}
	accountID, err := strconv.Atoi(c.Param("accountID"))
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "data": "账号ID格式错误"})
		return
	}

	// 删除相关规则
	rules, err := getRulesByAccountID(uint(accountID), userInfo.ID)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "获取相关规则失败"})
		return
	}
	for _, oneRule := range rules {
		err = rule.DeleteRuleByID(oneRule.ID, userInfo.ID)
		if err != nil {
			c.JSON(500, gin.H{"code": 500, "data": "删除相关规则失败"})
			return
		}
	}

	// 删除账号
	gormDB, err := db.GetGormDB()
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "服务器内部错误"})
		return
	}
	result := gormDB.Where("id = ? AND uid = ?", accountID, userInfo.ID).Delete(&model.Account{})
	if result.Error != nil {
		c.JSON(500, gin.H{"code": 500, "data": "删除账号失败"})
		return
	}
	c.JSON(200, gin.H{"code": 200, "data": "删除账号成功"})
}
