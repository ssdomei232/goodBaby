package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/ssdomei232/goodBaby/api/account"
	"github.com/ssdomei232/goodBaby/api/dashboard"
	apilog "github.com/ssdomei232/goodBaby/api/log"
	apimeta "github.com/ssdomei232/goodBaby/api/meta"
	"github.com/ssdomei232/goodBaby/api/rule"
	"github.com/ssdomei232/goodBaby/api/timer"
	"github.com/ssdomei232/goodBaby/api/user"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/handler/checker"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/handler/runner"
	"github.com/ssdomei232/goodBaby/web"
)

var (
	version   = "dev"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func main() {
	log.Printf("goodBaby %s (%s, %s) 启动中...", version, gitCommit, buildDate)

	config := configs.MustGetConfig()

	// 初始化数据库(建表/迁移)与执行器
	db.MustInit()
	runner.InitExecutorRegistry()

	// 启动定时检查
	c := cron.New()
	spec := fmt.Sprintf("@every %dm", config.CheckIntervalMinutes)
	if _, err := c.AddFunc(spec, checker.CheckTimers); err != nil {
		log.Fatalf("注册定时任务失败: %v", err)
	}
	c.Start()

	r := gin.Default()

	// 持久化的 session 密钥：重启后已登录用户不会掉线
	store := cookie.NewStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   config.SessionMaxAgeHours * 3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("goodbaby-session", store))

	// 前端本地开发时的 CORS 支持
	if len(config.AllowedOrigins) > 0 {
		r.Use(corsMiddleware(config.AllowedOrigins))
	}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/site", apimeta.HandleGetSiteInfo)
		v1.POST("/user/registry", user.HandleRegistry)
		v1.POST("/user/login", user.HandleLogin)
	}

	// 需要认证的路由组
	authorized := v1.Group("/")
	authorized.Use(user.AuthMiddleware())
	{
		authorized.GET("/providers", apimeta.HandleGetProviders)
		authorized.GET("/dashboard", dashboard.HandleGetOverview)

		users := authorized.Group("/user")
		{
			users.GET("/info", user.HandleGetUserInfo)
			users.POST("/logout", user.HandleLogout)
			users.POST("/password", user.HandleChangePassword)
			users.PUT("/notify", user.HandleUpdateNotifyConfig)
		}

		timers := authorized.Group("/timers")
		{
			timers.GET("/", timer.HandleGetAllTimers)
			timers.POST("/", timer.HandleCreateTimer)
			timers.POST("/sign", timer.HandleSignAll)
			timers.GET("/:timerID", timer.HandleGetTimer)
			timers.PUT("/:timerID", timer.HandleEditTimer)
			timers.POST("/:timerID/sign", timer.HandleSignTimer)
			timers.POST("/:timerID/trigger", timer.HandleTriggerTimer)
			timers.GET("/:timerID/check", timer.HandleCheckDeleteTimer)
			timers.DELETE("/:timerID", timer.HandleDeleteTimer)
		}

		rules := authorized.Group("/rules")
		{
			rules.GET("/", rule.HandleGetAllRules)
			rules.POST("/", rule.HandleCreateRule)
			rules.PUT("/:ruleID", rule.HandleEditRule)
			rules.POST("/:ruleID/test", rule.HandleTestRule)
			rules.DELETE("/:ruleID", rule.HandleDeleteRule)
		}

		accounts := authorized.Group("/accounts")
		{
			accounts.GET("/", account.HandleGetAllAccounts)
			accounts.POST("/", account.HandleAddAccount)
			accounts.PUT("/:accountID", account.HandleEditAccount)
			accounts.POST("/:accountID/test", account.HandleTestAccount)
			accounts.GET("/:accountID/check", account.HandleCheckDeleteAccount)
			accounts.DELETE("/:accountID", account.HandleDeleteAccount)
		}

		logs := authorized.Group("/logs")
		{
			logs.GET("/", apilog.HandleGetLogs)
			logs.DELETE("/", apilog.HandleClearLogs)
		}
	}

	// 挂载内嵌前端
	web.Register(r)

	log.Printf("goodBaby 正在监听 %s", config.ListenAddr)
	if err := r.Run(config.ListenAddr); err != nil {
		log.Fatalf("启动 HTTP 服务失败: %v", err)
	}
}

// corsMiddleware 允许配置中的来源跨域访问(携带 Cookie)
func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if _, ok := allowed[origin]; ok {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type")
			c.Header("Vary", "Origin")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
