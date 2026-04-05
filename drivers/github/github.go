package github

import (
	"context"
	"log"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/model"
)

// 将repos设置为public
func MakeRepositoryPublic(rule *model.Rule) {
	reposConfig, account := GetGithubReposAndAccountFromRule(rule)
	if reposConfig == nil || account == nil {
		log.Printf("无法获取 GitHub 配置或账户信息")
		return
	}

	for _, repo := range reposConfig.Repos {
		config, err := configs.GetConfig()
		if err != nil {
			log.Printf("获取配置失败: %v", err)
			return
		}
		var timeout time.Duration = time.Duration(config.TimeoutDurationHours) * time.Hour

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		operation := func() (string, error) {
			err := SetRepositoryPublic(account.Token, account.Owner, repo)
			return "", err
		}

		_, err = backoff.Retry(ctx, operation, backoff.WithBackOff(backoff.NewExponentialBackOff()))
		if err != nil {
			log.Printf("将仓库 %s 设置为 public 失败: %v", repo, err)
		}
	}
}
