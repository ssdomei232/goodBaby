package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/ssdomei232/goodBaby/configs"
)

func SendQQMsg() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Printf("获取配置文件失败: %v", err)
	}

	url := fmt.Sprintf("%s/%s", config.CatBotUrl, config.CatBotKey)

	for _, groupId := range config.QQSendGroup {
		sendMsg(url, groupId, config.QQMsg)
	}

}

func sendMsg(url string, groupId int, msg string) {
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("group_id", strconv.Itoa(groupId))
	_ = writer.WriteField("message", msg)
	err := writer.Close()
	if err != nil {
		log.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	var res *http.Response
	var maxRetries = 3
	for i := range maxRetries {
		res, err = client.Do(req)
		if err == nil {
			break // 请求成功，跳出循环
		}
		log.Printf("请求QQ API失败，正在重试... (%d/%d)", i+1, maxRetries)
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
}
