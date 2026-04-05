package account

import (
	"github.com/gin-gonic/gin"
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
