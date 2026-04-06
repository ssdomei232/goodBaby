package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/model"
)

// 处理注册请求，配置文件中可设置禁用注册
func HandleRegistry(c *gin.Context) {
	var userRegistryRequest model.UserRegistryReuest
	var err error

	config, err := configs.GetConfig()
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "获取配置文件失败"})
		return
	}

	if !config.EnableRegistry {
		c.JSON(403, gin.H{"code": 403, "data": "注册功能已关闭"})
		return
	}

	err = c.BindJSON(&userRegistryRequest)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "data": "输入错误"})
		return
	}

	user := model.User{
		Username: userRegistryRequest.Username,
		Password: userRegistryRequest.Password,
	}

	if err = user.IsValid(); err != nil {
		c.JSON(400, gin.H{"code": 400, "data": err.Error()})
		return
	}

	if isExist := user.IsExist(); isExist {
		c.JSON(400, gin.H{"code": 400, "data": "用户名已存在"})
		return
	}

	if err = createUser(&user); err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "注册失败"})
		return
	}

	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Save()

	c.JSON(200, gin.H{"code": 200, "data": "注册成功"})
}

// 处理登录请求
func HandleLogin(c *gin.Context) {
	var userLoginRequest model.UserRegistryReuest
	var err error

	err = c.BindJSON(&userLoginRequest)
	if err != nil {
		c.JSON(400, gin.H{"code": 400, "data": "输入错误"})
		return
	}

	user := model.User{
		Username: userLoginRequest.Username,
		Password: userLoginRequest.Password,
	}

	if err = user.IsValid(); err != nil {
		c.JSON(400, gin.H{"code": 400, "data": err.Error()})
	}

	if err = verifyUser(&user); err != nil {
		c.JSON(400, gin.H{"code": 400, "data": "用户名或密码错误"})
		return
	}

	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Save()

	c.JSON(200, gin.H{"code": 200, "data": "登录成功"})
}

// 处理获取用户信息请求
func HandleGetUserInfo(c *gin.Context) {
	userInfo, err := GetUserInfoByGinCtx(c)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "data": "获取用户信息失败"})
		return
	}
	c.JSON(200, gin.H{"code": 200, "data": userInfo})
}

// 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		if username == nil {
			c.JSON(401, gin.H{"code": 401, "data": "未登录"})
			c.Abort()
			return
		}
		c.Next()
	}
}
