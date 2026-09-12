package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v84/github"
)

// ParseAccountConfig 解析 GitHub 账号配置
func ParseAccountConfig(config string) (*GithubAccount, error) {
	if config == "" {
		return nil, fmt.Errorf("解析GitHub账号配置失败: 规则没有关联账号")
	}

	var account GithubAccount
	if err := json.Unmarshal([]byte(config), &account); err != nil {
		return nil, fmt.Errorf("解析GitHub账号配置失败: %v", err)
	}
	return &account, nil
}

// ParseReposConfig 解析仓库规则配置
func ParseReposConfig(configJSON string) (*GithubReposConfig, error) {
	var config GithubReposConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析GitHub规则配置失败: %v", err)
	}
	return &config, nil
}

// SetRepositoryPublic 把一个仓库设置为 public
func SetRepositoryPublic(ctx context.Context, token, owner, repo string) error {
	client := github.NewClient(nil).WithAuthToken(token)
	if _, _, err := client.Repositories.Edit(ctx, owner, repo, &github.Repository{
		Visibility: github.Ptr("public"),
	}); err != nil {
		return fmt.Errorf("设置仓库 %s 为 public 失败: %w", repo, err)
	}
	return nil
}
