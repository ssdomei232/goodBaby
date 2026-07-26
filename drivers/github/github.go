package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

// MakeRepositoryPublic 将规则中配置的仓库设置为 public
func MakeRepositoryPublic(ctx context.Context, rule *model.Rule) error {
	reposConfig, account, err := GetGithubReposAndAccountFromRule(rule)
	if err != nil {
		return err
	}

	var fails []string
	for _, repo := range reposConfig.Repos {
		err := retry.Do(ctx, func() error {
			return SetRepositoryPublic(ctx, account.Token, account.Owner, repo)
		})
		if err != nil {
			fails = append(fails, fmt.Sprintf("%s: %v", repo, err))
		}
	}

	if len(fails) > 0 {
		return fmt.Errorf("%d/%d 个仓库设置为 public 失败: %s",
			len(fails), len(reposConfig.Repos), strings.Join(fails, "; "))
	}
	return nil
}
