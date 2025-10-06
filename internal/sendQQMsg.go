package internal

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/ssdomei232/goodBaby/configs"
)

func SendQQ() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
		return // 添加return避免继续执行
	}

	url := fmt.Sprintf("%s/?secret=%s", config.CatBotUrl, config.CatBotKey)
	stopMsg := fmt.Sprintf("%s\n\n%s\n\n本消息由自动程序发送(摇篮系统)", GetBasicInfo(), config.QQMsg)

	for _, groupId := range config.QQSendGroup {
		sendQQMsg(url, groupId, stopMsg)
	}
}

func sendQQMsg(url string, groupId int, msg string) {
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	// 处理写入字段的错误
	if err := writer.WriteField("group_id", strconv.Itoa(groupId)); err != nil {
		log.Printf("写入group_id字段失败: %v", err)
		return
	}

	if err := writer.WriteField("message", msg); err != nil {
		log.Printf("写入message字段失败: %v", err)
		return
	}

	err := writer.Close()
	if err != nil {
		log.Printf("关闭multipart writer失败: %v", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Printf("创建HTTP请求失败: %v", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	var res *http.Response

	for i := range maxRetries {
		res, err = client.Do(req)
		if err == nil {
			break // 请求成功，跳出循环
		}
		if i < maxRetries {
			time.Sleep(retryDelay * time.Duration(i+1)) // 递增延迟
		}
		log.Printf("请求QQ API失败，正在重试... (%d/%d): %v", i+1, maxRetries, err)
	}

	// 检查res是否为nil，避免panic
	if res == nil {
		log.Println("QQ API请求失败，无法获取响应")
		return
	}

	// 使用defer确保资源被释放，但要确保res不为nil
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()

	// 可选：记录响应结果
	log.Println("QQ消息发送完成")
}
