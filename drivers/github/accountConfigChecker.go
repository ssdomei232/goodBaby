package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v84/github"
	"github.com/ssdomei232/goodBaby/internal/meta"
	"github.com/ssdomei232/goodBaby/internal/retry"
)

// GitHubAccountConfigValidator GitHub账号配置验证器
type GitHubAccountConfigValidator struct{}

func (v *GitHubAccountConfigValidator) GetType() string {
	return AccountType
}

func (v *GitHubAccountConfigValidator) Validate(config string) error {
	cfg, err := parseAccount(config)
	if err != nil {
		return err
	}

	if cfg.Token == "" {
		return fmt.Errorf("GitHub账号配置中Token不能为空")
	}
	if cfg.Owner == "" {
		return fmt.Errorf("GitHub账号配置中Owner不能为空")
	}

	return nil
}

// Test 用 token 拉一次当前用户信息，验证 token 是否有效
func (v *GitHubAccountConfigValidator) Test(config string) error {
	cfg, err := parseAccount(config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), retry.TestTimeout)
	defer cancel()

	client := github.NewClient(nil).WithAuthToken(cfg.Token)
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("GitHub Token 校验失败: %w", err)
	}
	if user == nil || user.GetLogin() == "" {
		return fmt.Errorf("GitHub Token 无效")
	}
	return nil
}

func (v *GitHubAccountConfigValidator) Meta() meta.AccountMeta {
	return meta.AccountMeta{
		Type:        AccountType,
		Label:       "GitHub",
		Description: "使用 Personal Access Token 操作仓库，需要 repo 权限。",
		Fields: []meta.Field{
			{Key: "token", Label: "Personal Access Token", Type: meta.FieldPassword, Required: true, Secret: true, Placeholder: "ghp_xxx", Help: "需要勾选 repo 权限"},
			{Key: "owner", Label: "仓库所有者", Type: meta.FieldString, Required: true, Placeholder: "your-github-name", Help: "用户名或组织名"},
		},
	}
}

func parseAccount(config string) (*GithubAccount, error) {
	var cfg GithubAccount
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return nil, fmt.Errorf("解析GitHub账号配置失败: %v", err)
	}
	return &cfg, nil
}
