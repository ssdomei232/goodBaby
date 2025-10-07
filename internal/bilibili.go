package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/ssdomei232/goodBaby/configs"
)

func isCookiesExpired(cookiesString string) bool {
	// 解析cookie中的过期时间
	lines := strings.SplitSeq(cookiesString, "\n")
	for line := range lines {
		if strings.Contains(line, "Expires=") {
			// 提取过期时间字符串
			expireStr := ""
			parts := strings.Split(line, ";")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "Expires=") {
					expireStr = strings.TrimPrefix(part, "Expires=")
					break
				}
			}

			if expireStr != "" {
				// 解析过期时间
				expireTime, err := http.ParseTime(expireStr)
				if err != nil {
					log.Printf("解析过期时间失败: %v", err)
					return true // 如果解析失败，认为已过期
				}

				// 比较过期时间和当前时间
				if time.Now().After(expireTime) {
					return true // 已过期
				}
			}
		}
	}

	return false // 未过期
}

// 二维码登录
func LoginWithQRCode(biliClient *bilibili.Client) {
	qrCode, err := biliClient.GetQRCode()
	if err != nil {
		log.Printf("获取二维码失败: %v", err)
		return
	}

	log.Println("请使用哔哩哔哩APP扫码登录或进入文件目录寻找qrcode.png查看二维码")
	buf, _ := qrCode.Encode()
	os.WriteFile("qrcode.png", buf, 0644)
	qrCode.Print()

	result, err := biliClient.LoginWithQRCode(bilibili.LoginWithQRCodeParam{
		QrcodeKey: qrCode.QrcodeKey,
	})
	if err != nil || result.Code != 0 {
		log.Printf("登录失败: %v", err)
		return
	}

	log.Println("登录成功")

	// 保存新获取的cookie
	saveCookies(biliClient)
}

// 保存cookie到文件
func saveCookies(client *bilibili.Client) {
	cookiesString := client.GetCookiesString()
	err := os.WriteFile("cookies.txt", []byte(cookiesString), 0644)
	if err != nil {
		log.Printf("保存cookie失败: %v", err)
	}
}

// 从文件加载cookie并验证有效性
func LoadCookies(client *bilibili.Client) bool {
	// 检查cookie文件是否存在
	if _, err := os.Stat("cookies.txt"); os.IsNotExist(err) {
		return false
	}

	// 读取cookie
	cookiesBytes, err := os.ReadFile("cookies.txt")
	if err != nil {
		log.Printf("读取cookie文件失败: %v", err)
		return false
	}

	cookiesString := string(cookiesBytes)
	client.SetCookiesString(cookiesString)

	// 验证cookie是否过期
	if isCookiesExpired(cookiesString) {
		log.Println("cookie已过期")
		return false
	}

	return true
}

// 在Bilibili动态发送死亡通告
func SendBili(biliClient *bilibili.Client) {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
		return
	}

	stopMsg := fmt.Sprintf("%s\n\n%s\n\n本消息由自动程序发送(摇篮系统)", GetBasicInfo(), config.BiliMsg)

	var dynamicParams bilibili.CreateDynamicParam
	dynamicParams = bilibili.CreateDynamicParam{
		DynamicId: 0,
		Type:      4,
		Rid:       0,
		Content:   stopMsg,
	}

	// 重试机制
	retryCount := 10         // 重试次数
	delay := time.Second * 2 // 初始重试延迟

	for i := 0; i <= retryCount; i++ {
		_, err := biliClient.CreateDynamic(dynamicParams)
		if err == nil {
			log.Println("动态发送成功")
			return
		}

		log.Printf("动态发送失败 (尝试 %d/%d): %v", i+1, retryCount+1, err)

		// 如果不是最后一次尝试，则等待后重试
		if i < retryCount {
			log.Printf("等待 %v 后进行第 %d 次重试", delay, i+2)
			time.Sleep(delay)
			delay *= 2 // 指数退避
		}
	}

	log.Println("动态发送最终失败，已达到最大重试次数")
}

// 检查cookie有效性
func CheckCookieValidity(biliClient *bilibili.Client) {
	// 读取存储的cookie文件
	if _, err := os.Stat("cookies.txt"); os.IsNotExist(err) {
		log.Println("cookie文件不存在,需要重新登录")
		triggerLoginRequest(biliClient)
		return
	}

	cookiesBytes, err := os.ReadFile("cookies.txt")
	if err != nil {
		log.Printf("读取cookie文件失败: %v", err)
		triggerLoginRequest(biliClient)
		return
	}

	cookiesString := string(cookiesBytes)

	// 检查cookie是否即将过期（提前1天提醒）
	if isCookiesExpiringSoon(cookiesString) {
		log.Println("cookie即将过期,需要重新登录")
		triggerLoginRequest(biliClient)
	}
}

// 检查cookie是否即将过期（提前1天）
func isCookiesExpiringSoon(cookiesString string) bool {
	lines := strings.Split(cookiesString, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Expires=") {
			// 提取过期时间字符串
			expireStr := ""
			parts := strings.Split(line, ";")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "Expires=") {
					expireStr = strings.TrimPrefix(part, "Expires=")
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
	// 获取配置
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
		return
	}

	// 生成新的二维码
	qrCode, err := biliClient.GetQRCode()
	if err != nil {
		log.Printf("获取二维码失败: %v", err)
		return
	}

	// 在命令行打印二维码
	log.Println("需要重新登录哔哩哔哩")
	qrCode.Print()

	// 通过邮件发送登录请求
	emailTitle := "B站登录请求 - Cookie即将过期"
	emailMsg := "您的B站cookie即将过期,请尽快重启摇篮系统进行登陆,请尽快完成登录操作。"

	// 发送邮件通知
	if config.BiliWarnAddress != "" {
		err = sendMailMsg(config.BiliWarnAddress, emailMsg, emailTitle)
		if err != nil {
			log.Printf("发送登录邮件失败: %v", err)
		} else {
			log.Println("登录请求邮件已发送")
		}
	}
}
