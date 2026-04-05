package github

type GithubReposConfig struct {
	Repos []string `yaml:"repos"`
}

type GithubAccount struct {
	Token string
	Owner string
}
