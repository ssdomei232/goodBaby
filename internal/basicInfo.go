package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/ssdomei232/goodBaby/configs"
)

func GetBasicInfo() string {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
		return "获取配置文件失败"
	}

	now := time.Now()
	lastSignalTime := now.Add(-time.Duration(config.DisconnectDuration) * time.Hour).Format("2006-01-02 15:04:05")

	if config.Debug {
		return fmt.Sprintf("此消息为测试消息\n\n%s(QQ号:%s)已于%s至%s(UTC+8)之间去世,享年%d岁,死因可能为%s(此死因为推测最可能的死因,请勿完全相信,请等待后续消息)",
			config.Basic.Name, config.Basic.QQNumber, lastSignalTime, time.Now().Format("2006-01-02 15:04:05"), config.Basic.Age, config.Basic.CauseStop)
	} else {
		return fmt.Sprintf("%s(QQ号:%s,网名:%s)已于%s至%s(UTC+8)之间去世,享年%d岁,死因可能为%s(此死因为推测最可能的死因,请勿完全相信,请等待后续消息)",
			config.Basic.Name, config.Basic.NickName, config.Basic.QQNumber, lastSignalTime, time.Now().Format("2006-01-02 15:04:05"), config.Basic.Age, config.Basic.CauseStop)
	}

}
