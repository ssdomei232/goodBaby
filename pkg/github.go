package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// RepositoryVisibility 用于更新仓库可见性的请求体
type RepositoryVisibility struct {
	Private bool `json:"private"`
}

// MakeRepositoryPublic 将指定的GitHub仓库从私有设置为公开
func MakeRepositoryPublic(owner, repo, token string) error {
	// GitHub API URL for updating repository visibility
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	// 创建请求体，将 private 设置为 false 表示公开仓库
	visibility := RepositoryVisibility{
		Private: false, // false 表示公开仓库
	}

	// 将请求体序列化为 JSON
	jsonData, err := json.Marshal(visibility)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 创建 PATCH 请求
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	log.Printf("仓库 %s/%s 已成功设置为公开", owner, repo)
	return nil
}

// MakeRepositoryPrivate 将指定的GitHub仓库从公开设置为私有（附加功能）
func MakeRepositoryPrivate(owner, repo, token string) error {
	// GitHub API URL for updating repository visibility
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	// 创建请求体，将 private 设置为 true 表示私有仓库
	visibility := RepositoryVisibility{
		Private: true, // true 表示私有仓库
	}

	// 将请求体序列化为 JSON
	jsonData, err := json.Marshal(visibility)
	if err != nil {
		return fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 创建 PATCH 请求
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	log.Printf("仓库 %s/%s 已成功设置为私有", owner, repo)
	return nil
}
