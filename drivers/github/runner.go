package github

import (
	"context"

	"github.com/ssdomei232/goodBaby/model"
)

// GithubMakeRepoPublicExecutor GitHub仓库公开执行器
type GithubMakeRepoPublicExecutor struct{}

func (e *GithubMakeRepoPublicExecutor) GetType() string {
	return RuleTypeMakeRepoPublic
}

func (e *GithubMakeRepoPublicExecutor) Execute(ctx context.Context, task *model.RuleTask) error {
	account, err := ParseAccountConfig(task.AccountConfig())
	if err != nil {
		return err
	}

	config, err := ParseReposConfig(task.Rule.ConfigJson)
	if err != nil {
		return err
	}

	return MakeRepositoryPublic(ctx, account, config)
}
