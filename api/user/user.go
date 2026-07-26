package user

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/response"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// HandleRegistry 处理注册请求，配置文件中可设置禁用注册
//
// 系统内还没有任何用户时始终允许注册，否则全新部署将无法创建第一个账号。
func HandleRegistry(c *gin.Context) {
	config, err := configs.GetConfig()
	if err != nil {
		response.ServerError(c, "获取配置文件失败")
		return
	}

	userCount, err := CountUsers()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	if !config.EnableRegistry && userCount > 0 {
		response.Fail(c, http.StatusForbidden, "注册功能已关闭")
		return
	}

	var registryRequest model.UserRegistryReuest
	if err := c.ShouldBindJSON(&registryRequest); err != nil {
		response.BadRequest(c, "输入错误")
		return
	}

	user := model.User{
		Username: registryRequest.Username,
		Password: registryRequest.Password,
	}

	if err := user.IsValid(); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	taken, err := IsUsernameTaken(user.Username)
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}
	if taken {
		response.BadRequest(c, "用户名已存在")
		return
	}

	if err := createUser(&user); err != nil {
		response.ServerError(c, "注册失败")
		return
	}

	if err := setSession(c, &user); err != nil {
		response.ServerError(c, "写入会话失败")
		return
	}

	response.OK(c, "注册成功")
}

// HandleLogin 处理登录请求
func HandleLogin(c *gin.Context) {
	var loginRequest model.UserRegistryReuest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		response.BadRequest(c, "输入错误")
		return
	}

	if loginRequest.Username == "" || loginRequest.Password == "" {
		response.BadRequest(c, "用户名或密码不能为空")
		return
	}

	user, err := verifyUser(loginRequest.Username, loginRequest.Password)
	switch {
	case errors.Is(err, ErrUserNotFound), errors.Is(err, ErrWrongPassword):
		// 不区分“用户不存在”和“密码错误”，避免泄露用户名是否存在
		response.Unauthorized(c, "用户名或密码错误")
		return
	case err != nil:
		response.ServerError(c, "登录失败")
		return
	}

	if err := setSession(c, user); err != nil {
		response.ServerError(c, "写入会话失败")
		return
	}

	response.OK(c, "登录成功")
}

// HandleLogout 退出登录
func HandleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	if err := session.Save(); err != nil {
		response.ServerError(c, "退出登录失败")
		return
	}
	response.OK(c, "已退出登录")
}

// HandleGetUserInfo 获取当前登录用户信息
func HandleGetUserInfo(c *gin.Context) {
	userInfo, err := GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}
	response.OK(c, userInfo)
}

// HandleChangePassword 修改密码
func HandleChangePassword(c *gin.Context) {
	userInfo, err := GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入错误")
		return
	}

	if _, err := verifyUser(userInfo.Username, req.OldPassword); err != nil {
		response.BadRequest(c, "原密码错误")
		return
	}

	candidate := model.User{Username: userInfo.Username, Password: req.NewPassword}
	if err := candidate.IsValid(); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	hashed, err := encryptPassword(req.NewPassword)
	if err != nil {
		response.ServerError(c, "修改密码失败")
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}
	if err := gormDB.Model(&model.User{}).Where("id = ?", userInfo.ID).
		Update("password", hashed).Error; err != nil {
		response.ServerError(c, "修改密码失败")
		return
	}

	response.OK(c, "密码修改成功")
}

// HandleUpdateNotifyConfig 更新提醒渠道(钉钉机器人)配置
func HandleUpdateNotifyConfig(c *gin.Context) {
	userInfo, err := GetUserInfoByGinCtx(c)
	if err != nil {
		response.Unauthorized(c, "获取用户信息失败")
		return
	}

	var req model.NotifyConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "输入错误")
		return
	}

	if err := validateNotifyConfig(req.DingTalkConfig); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	gormDB, err := db.GetGormDB()
	if err != nil {
		response.ServerError(c, "服务器内部错误")
		return
	}

	if err := gormDB.Model(&model.User{}).Where("id = ?", userInfo.ID).
		Update("ding_talk_config", req.DingTalkConfig).Error; err != nil {
		response.ServerError(c, "保存提醒配置失败")
		return
	}

	response.OK(c, "提醒配置已保存")
}

// AuthMiddleware 认证中间件，同时把当前用户放进 context 供后续 handler 复用
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		uid, ok := session.Get("uid").(uint)
		if !ok {
			response.AbortWith(c, http.StatusUnauthorized, "未登录")
			return
		}

		user, err := GetUserByID(uid)
		if err != nil {
			// 用户已被删除，清理会话
			session.Clear()
			_ = session.Save()
			response.AbortWith(c, http.StatusUnauthorized, "未登录")
			return
		}

		c.Set(contextKey, user)
		c.Next()
	}
}
