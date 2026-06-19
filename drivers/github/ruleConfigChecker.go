package github

import (
	"encoding/json"
	"fmt"
)

type GithubMakeRepositoryPublicRuleValidator struct{}

func (v *GithubMakeRepositoryPublicRuleValidator) GetType() string {
	return "github_make_repository_public"
}

func (v *GithubMakeRepositoryPublicRuleValidator) Validate(configJSON string) error {
	var config GithubReposConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析 GitHub 规则配置失败: %v", err)
	}

	if len(config.Repos) == 0 {
		return fmt.Errorf("GitHub 规则配置错误: 至少需要指定一个仓库")
	}

	return nil
}
