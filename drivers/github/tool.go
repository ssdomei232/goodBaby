package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v84/github"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// GetGithubReposAndAccountFromRule 从 rule 中获取 GithubReposConfig 和 GithubAccount
func GetGithubReposAndAccountFromRule(rule *model.Rule) (*GithubReposConfig, *GithubAccount, error) {
	var githubReposConfig GithubReposConfig
	if err := json.Unmarshal([]byte(rule.ConfigJson), &githubReposConfig); err != nil {
		return nil, nil, fmt.Errorf("解析 GitHub 规则配置失败: %w", err)
	}

	var githubAccount GithubAccount
	if err := db.LoadAccountConfig(rule.AccountID, &githubAccount); err != nil {
		return nil, nil, err
	}

	return &githubReposConfig, &githubAccount, nil
}

// SetRepositoryPublic 将仓库设置为 public
func SetRepositoryPublic(ctx context.Context, token, owner, repo string) error {
	client := github.NewClient(nil).WithAuthToken(token)

	opts := &github.Repository{
		Visibility: github.Ptr("public"),
	}

	_, _, err := client.Repositories.Edit(ctx, owner, repo, opts)
	return err
}
