package github

const (
	// AccountType GitHub 账号类型标识
	AccountType = "github"
	// RuleTypeMakeRepoPublic 公开仓库的规则类型标识
	RuleTypeMakeRepoPublic = "github-repo-public"
)

type GithubReposConfig struct {
	Repos []string `json:"repos"`
}

type GithubAccount struct {
	Token string `json:"token"`
	Owner string `json:"owner"`
}
