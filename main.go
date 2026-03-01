package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/internal/frontend"
	"github.com/ssdomei232/goodBaby/internal/reminder"
	"github.com/ssdomei232/goodBaby/internal/sender"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

var timer *time.Timer
var duration time.Duration
var biliClient *bilibili.Client
var cookieCheckTimer *time.Ticker

func init() {
	// 1. 加载配置文件
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
	}

	if config.Debug {
		duration = time.Duration(config.DisconnectDuration) * time.Second
	} else {
		duration = time.Duration(config.DisconnectDuration) * time.Hour
	}

	// 2. 检查并创建tmp目录
	if err := sender.EnsureTmpDirectory(); err != nil {
		log.Printf("创建tmp目录失败: %v", err)
	}

	// 3. 初始化哔哩哔哩和定时器
	biliClient = bilibili.New()
	sender.InitBilibili(biliClient)
	sender.StartCookieChecker(cookieCheckTimer, biliClient) // Check Bilibili cookie

	reminder.InitTimerManager(duration)
	timer = time.NewTimer(duration)
	go func() {
		<-timer.C
		trigger(config)
	}()
}

// 触发器
func trigger(config configs.Config) {
	go sender.SendOneBot()
	go sender.SendMail()
	go sender.Github()
	if !config.Debug {
		go sender.SendBili(biliClient)
	}
}

func main() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
	}
	// 定时任务
	reminder.Reminder()
	c := cron.New()
	if config.Debug {
		c.AddFunc("@every 1s", reminder.Reminder)
	} else {
		c.AddFunc("@every 1h", reminder.Reminder)
	}
	c.Start()

	// 前端/后端
	r := gin.Default()
	templFS, _ := fs.Sub(templateFiles, "templates")
	r.SetHTMLTemplate(frontend.LoadTemplates(templFS))
	staticFS, _ := fs.Sub(staticFiles, "static")
	r.StaticFS("/static", http.FS(staticFS))

	r.GET("/", frontend.IndexPage)
	r.GET("/signal", handleSignal)
	r.GET("/timer/status", reminder.HandleTimerStatus)
	r.Run(":8088")
}

// Handle Signal
func handleSignal(c *gin.Context) {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
	}

	secret := c.Query("secret")
	if secret != config.SignalSecret {
		c.JSON(403, gin.H{
			"code":    403,
			"message": "secret error",
		})
		return
	}
	timer.Reset(duration)
	reminder.GlobalTimerManager.Reset()
	c.JSON(200, gin.H{
		"code":    200,
		"message": "ok",
	})
	log.Println("触发信号")
}
