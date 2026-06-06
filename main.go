package main

import (
	"crypto/rand"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/ssdomei232/goodBaby/api/account"
	"github.com/ssdomei232/goodBaby/api/rule"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/handler/checker"
	"github.com/ssdomei232/goodBaby/handler/runner"
)

func main() {
	runner.InitExecutorRegistry()

	c := cron.New()
	c.AddFunc("@every 10m", checker.CheckTimers)
	c.Start()

	r := gin.Default()
	store := cookie.NewStore(generateRandomKey(32))
	r.Use(sessions.Sessions("goodbaby-session", store))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/user/registry", user.HandleRegistry)
		v1.POST("/user/login", user.HandleLogin)
	}

	// 需要认证的路由组
	authorized := v1.Group("/")
	authorized.Use(user.AuthMiddleware())
	{
		authorized.GET("/user/info", user.HandleGetUserInfo)
		authorized.GET("/rules", rule.HandleGetAllRules)
		authorized.POST("/rules", rule.HandleCreateRule)
		authorized.DELETE("/rules/:ruleID", rule.HandleDeleteRule)
		authorized.GET("/accounts", account.HandleGetAllAccounts)
		authorized.POST("/accounts", account.HandleAddAccount)
		authorized.GET("/accounts/:accountID/check", account.HandleCheckDeleteAccount)
		authorized.DELETE("/accounts/:accountID", account.HandleDeleteAccount)
	}

	r.Run(":8088")
}

func generateRandomKey(length int) []byte {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		panic(fmt.Sprintf("无法生成随机密钥: %v", err))
	}
	return key
}
