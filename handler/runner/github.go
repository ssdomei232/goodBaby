package runner

import (
	"log"

	"github.com/ssdomei232/goodBaby/drivers/github"
	"github.com/ssdomei232/goodBaby/model"
)

// GithubMakeRepoPublicExecutor GitHub仓库公开执行器
type GithubMakeRepoPublicExecutor struct{}

func (e *GithubMakeRepoPublicExecutor) GetType() string {
	return "github-make-repo-public"
}

func (e *GithubMakeRepoPublicExecutor) Execute(rule *model.Rule) error {
	log.Printf("执行GitHub仓库公开规则: %s (ID: %d)", rule.Name, rule.ID)

	github.MakeRepositoryPublic(rule)

	return nil
}
