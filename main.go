// main.go
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
	"github.com/ssdomei232/goodBaby/internal"
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
	// Get COnfig
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
	}

	if config.Debug {
		duration = time.Duration(config.DisconnectDuration) * time.Second
	} else {
		duration = time.Duration(config.DisconnectDuration) * time.Hour
	}

	// Check and create tmp dir
	if err := internal.EnsureTmpDirectory(); err != nil {
		log.Printf("创建tmp目录失败: %v", err)
	}

	// Init Bilibili Clinet
	biliClient = bilibili.New()
	internal.InitBilibili(biliClient)
	internal.StartCookieChecker(cookieCheckTimer, biliClient) // Check Bilibili cookie

	internal.InitTimerManager(duration)
	timer = time.NewTimer(duration)
	go func() {
		<-timer.C
		trigger(config)
	}()
}

// trigger
func trigger(config configs.Config) {
	if config.EnableQQ {
		go internal.SendQQ()
	}
	go internal.SendMail()
	go internal.Github()
	if !config.Debug {
		go internal.SendBili(biliClient)
	}
}

func main() {
	// cron job
	internal.Reminder()
	c := cron.New()
	c.AddFunc("@every 1h", internal.Reminder)
	c.Start()

	// Web service
	r := gin.Default()
	templFS, _ := fs.Sub(templateFiles, "templates")
	r.SetHTMLTemplate(internal.LoadTemplates(templFS))
	staticFS, _ := fs.Sub(staticFiles, "static")
	r.StaticFS("/static", http.FS(staticFS))

	r.GET("/", internal.IndexPage)
	r.GET("/signal", handleSignal)
	r.GET("/timer/status", internal.HandleTimerStatus)
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
	internal.GlobalTimerManager.Reset()
	c.JSON(200, gin.H{
		"code":    200,
		"message": "ok",
	})
	log.Println("触发信号")
}
