// main.go
package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/gin-gonic/gin"
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
	initBilibili()
	startCookieChecker() // Check Bilibili cookie

	internal.InitTimerManager(duration)
	timer = time.NewTimer(duration)
	go func() {
		<-timer.C
		trigger(config)
	}()
}

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
	r := gin.Default()
	templFS, _ := fs.Sub(templateFiles, "templates")
	r.SetHTMLTemplate(loadTemplates(templFS))
	staticFS, _ := fs.Sub(staticFiles, "static")
	r.StaticFS("/static", http.FS(staticFS))

	r.GET("/", indexPage)
	r.GET("/signal", handleSignal)
	r.GET("/timer/status", internal.HandleTimerStatus)
	r.Run(":8088")
}

func indexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func initBilibili() {
	// try cached cookie
	if !internal.LoadCookies(biliClient) {
		// if cookie check failed, request qrcode login
		internal.LoginWithQRCode(biliClient)
	} else {
		log.Println("使用已存储的有效cookie登录")
	}
}

// Enable periodic cookie checks
func startCookieChecker() {
	// check cookie per hour
	cookieCheckTimer = time.NewTicker(1 * time.Hour)
	go func() {
		for {
			<-cookieCheckTimer.C
			internal.CheckCookieValidity(biliClient)
		}
	}()
}

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

// loadTemplates loads templates from the embedded file system
func loadTemplates(filesystem fs.FS) *template.Template {
	templ := template.Must(template.New("").ParseFS(filesystem, "*.html"))
	return templ
}
