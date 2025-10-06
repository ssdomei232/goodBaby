package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/internal"
)

var timer *time.Timer
var duration time.Duration

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
}

func trigger() {
	go internal.SendQQ()
	go internal.SendMail()
}

func main() {
	r := gin.Default()
	r.GET("/signal", func(c *gin.Context) {
		timer.Reset(duration)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "ok",
		})
	})
	r.Run(":8088")
}
