package main

import (
	"log"
	"time"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/internal"
)

var timer *time.Timer
var duration time.Duration
var biliClient *bilibili.Client
var cookieCheckTimer *time.Ticker

func init() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
	}

	if config.Debug {
		duration = time.Duration(config.DisconnectDuration) * time.Second
	} else {
		duration = time.Duration(config.DisconnectDuration) * time.Hour
	}

	timer = time.NewTimer(duration)
	go func() {
		<-timer.C
		trigger()
	}()

	// 初始化bilibili客户端
	biliClient = bilibili.New()
	initBilibili()
	startCookieChecker() // 启动定期检查cookie有效性
}

func trigger() {
	go internal.SendQQ()
	go internal.SendMail()
	go internal.Github()
	// go internal.SendBili(biliClient)
}

func main() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
	}

	r := gin.Default()
	r.GET("/signal", func(c *gin.Context) {
		secret := c.Query("secret")
		if secret != config.SignalSecret {
			c.JSON(403, gin.H{
				"code":    403,
				"message": "secret error",
			})
			return
		}
		timer.Reset(duration)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "ok",
		})
		log.Println("触发信号")
	})
	r.Run(":8088")
}

func initBilibili() {
	// 尝试加载已存储的cookie
	if !internal.LoadCookies(biliClient) {
		// 如果没有有效cookie，则进行二维码登录
		internal.LoginWithQRCode(biliClient)
	} else {
		log.Println("使用已存储的有效cookie登录")
	}

}

// 启动cookie定期检查
func startCookieChecker() {
	// 每小时检查一次cookie有效性
	cookieCheckTimer = time.NewTicker(1 * time.Hour)
	go func() {
		for {
			<-cookieCheckTimer.C
			internal.CheckCookieValidity(biliClient)
		}
	}()
}
