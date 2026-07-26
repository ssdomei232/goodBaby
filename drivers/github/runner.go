package github

import (
	"context"
	"log"

	"github.com/ssdomei232/goodBaby/model"
)

// GithubMakeRepoPublicExecutor GitHub仓库公开执行器
type GithubMakeRepoPublicExecutor struct{}

func (e *GithubMakeRepoPublicExecutor) GetType() string {
	return RuleTypeMakeRepoPublic
}

func (e *GithubMakeRepoPublicExecutor) Execute(ctx context.Context, rule *model.Rule) error {
	log.Printf("执行GitHub仓库公开规则: %s (ID: %d)", rule.Name, rule.ID)

	return MakeRepositoryPublic(ctx, rule)
}
