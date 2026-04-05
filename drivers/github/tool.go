package github

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/go-github/v84/github"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// 从 rule 中获取 GithubReposConfig 和 GithubAccount
func GetGithubReposAndAccountFromRule(rule *model.Rule) (*GithubReposConfig, *GithubAccount) {
	var githubReposConfig GithubReposConfig
	var githubAccount GithubAccount

	// 获取 GithubReposConfig
	err := json.Unmarshal([]byte(rule.ConfigJson), &githubReposConfig)
	if err != nil {
		log.Printf("获取 GitHub 配置失败: %v", err)
		return nil, nil
	}

	// 获取 GithubAccount
	gormDB, err := db.GetGormDB()
	if err != nil {
		log.Printf("获取 Gorm DB 失败: %v", err)
		return nil, nil
	}

	err = gormDB.Where("id = ?", rule.AccountID).First(&githubAccount).Error
	if err != nil {
		log.Printf("获取 GitHub 账户失败: %v", err)
		return nil, nil
	}

	return &githubReposConfig, &githubAccount
}

// 将仓库设置为 public
func SetRepositoryPublic(token, owner, repo string) error {
	ctx := context.Background()

	// 1. 初始化客户端
	client := github.NewClient(nil).WithAuthToken(token)

	// 2. 准备修改的参数
	opts := &github.Repository{
		Visibility: github.Ptr("public"),
	}

	// 3. 执行更新操作
	_, _, err := client.Repositories.Edit(ctx, owner, repo, opts)

	return err
}
