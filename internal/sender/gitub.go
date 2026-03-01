package sender

import (
	"log"
	"time"

	"github.com/ssdomei232/goodBaby/configs"
	"github.com/ssdomei232/goodBaby/pkg/github"
)

// 将repos设置为public
func Github() {
	config, err := configs.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	for _, repo := range config.GithubConfig.Repos {
		maxRetries := 10
		for attempt := range maxRetries {
			err := github.MakeRepositoryPublic(config.GithubConfig.Owner, repo, config.GithubConfig.Token)
			if err == nil {
				break
			}

			log.Printf("尝试设置仓库 %s 为公开失败 (尝试 %d/%d): %v", repo, attempt+1, maxRetries, err)

			// 如果不是最后一次尝试，则等待一段时间后重试
			if attempt < maxRetries-1 {
				time.Sleep(time.Duration(attempt+1) * time.Second) // 逐步增加等待时间
			} else {
				log.Printf("设置仓库 %s 为公开最终失败", repo)
			}
		}
	}
}
