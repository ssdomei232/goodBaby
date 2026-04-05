package main

import (
	"crypto/rand"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/api/account"
	"github.com/ssdomei232/goodBaby/api/rule"
	"github.com/ssdomei232/goodBaby/api/user"
)

func main() {
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
