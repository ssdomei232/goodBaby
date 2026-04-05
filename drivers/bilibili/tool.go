package bilibili

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// 检查cookie是否过期
func isCookiesExpired(cookiesString string) bool {
	// 从cookie字符串中提取过期时间
	lines := strings.SplitSeq(cookiesString, "\n")
	for line := range lines {
		if strings.Contains(line, "Expires=") {
			// 提取过期时间字符串
			expireStr := ""
			parts := strings.SplitSeq(line, ";")
			for part := range parts {
				part = strings.TrimSpace(part)
				if after, ok := strings.CutPrefix(part, "Expires="); ok {
					expireStr = after
					break
				}
			}

			if expireStr != "" {
				// 解析过期时间
				expireTime, err := http.ParseTime(expireStr)
				if err != nil {
					log.Printf("解析过期时间失败: %v", err)
					return true // 如果解析失败，cookie已过期
				}
				if time.Now().After(expireTime) {
					return true
				}
			}
		}
	}

	return false
}

// 检查cookie是否即将过期（提前1天）
func isCookiesExpiringSoon(cookiesString string) bool {
	lines := strings.SplitSeq(cookiesString, "\n")
	for line := range lines {
		if strings.Contains(line, "Expires=") {
			// 提取过期时间字符串
			expireStr := ""
			parts := strings.SplitSeq(line, ";")
			for part := range parts {
				part = strings.TrimSpace(part)
				if after, ok := strings.CutPrefix(part, "Expires="); ok {
					expireStr = after
					break
				}
			}

			if expireStr != "" {
				// 解析过期时间
				expireTime, err := http.ParseTime(expireStr)
				if err != nil {
					log.Printf("解析过期时间失败: %v", err)
					return true // 如果解析失败，认为需要重新登录
				}

				// 检查是否即将过期（提前1天）
				if time.Now().Add(24 * time.Hour).After(expireTime) {
					return true // 即将过期
				}
			}
		}
	}

	return false
}

// 触发登录请求
func triggerLoginRequest(biliClient *bilibili.Client) {

	// 生成新的二维码
	qrCode, err := biliClient.GetQRCode()
	if err != nil {
		log.Printf("获取二维码失败: %v", err)
		return
	}

	// 在命令行打印二维码
	log.Println("需要重新登录哔哩哔哩")
	qrCode.Print()

	// 通过Dingtalk发送登录请求
	// title := "B站登录请求 - Cookie即将过期"
	// msg := "您的B站cookie即将过期,请尽快重启摇篮系统进行登陆,请尽快完成登录操作。"

	// 发送Dingtalk通知
	// if config.DingtalkBot.AccessToken != "" && config.DingtalkBot.Secret != "" {
	// 	SendDingTalkMsg(title, msg, config.DingtalkBot.AccessToken, config.DingtalkBot.Secret)
	// }

}

// 获取 Bilibili Client
func GetBiliClient(rule *model.Rule) (*bilibili.Client, error) {
	client := bilibili.New()
	gormDB, err := db.GetGormDB()
	if err != nil {
		return nil, err
	}

	// 1. 通过 account_id 获取到对应的 account 配置
	var biliAccount model.Account
	result := gormDB.Where("id = ?", rule.AccountID).First(&biliAccount)
	if result.Error != nil {
		return nil, result.Error
	}

	// 2. 通过 account 配置中的 config 字段获取到 CookiesString
	var biliAccountConfig BiliAccount
	err = json.Unmarshal([]byte(biliAccount.Config), &biliAccountConfig)
	if err != nil {
		return nil, err
	}

	client.SetCookiesString(biliAccountConfig.CookiesString)

	return client, nil
}
