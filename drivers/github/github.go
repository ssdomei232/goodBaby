package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/ssdomei232/goodBaby/internal/retry"
)

// MakeRepositoryPublic 把配置里的仓库逐个设置为 public
//
// 每个仓库单独重试、互不阻塞，最后汇总失败信息交给上层记录。
func MakeRepositoryPublic(ctx context.Context, account *GithubAccount, config *GithubReposConfig) error {
	var fails []string
	for _, repo := range config.Repos {
		if err := retry.Do(ctx, func() error {
			return SetRepositoryPublic(ctx, account.Token, account.Owner, repo)
		}); err != nil {
			fails = append(fails, err.Error())
		}
	}

	if len(fails) > 0 {
		return fmt.Errorf("%d/%d 个仓库设置为 public 失败: %s",
			len(fails), len(config.Repos), strings.Join(fails, "; "))
	}
	return nil
}
