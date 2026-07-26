package github

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ssdomei232/goodBaby/internal/meta"
)

type GithubMakeRepositoryPublicRuleValidator struct{}

func (v *GithubMakeRepositoryPublicRuleValidator) GetType() string {
	return RuleTypeMakeRepoPublic
}

func (v *GithubMakeRepositoryPublicRuleValidator) Validate(configJSON string) error {
	var config GithubReposConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return fmt.Errorf("解析 GitHub 规则配置失败: %v", err)
	}

	if len(config.Repos) == 0 {
		return fmt.Errorf("GitHub 规则配置错误: 至少需要指定一个仓库")
	}

	for _, repo := range config.Repos {
		if strings.TrimSpace(repo) == "" {
			return fmt.Errorf("GitHub 规则配置错误: 仓库名不能为空")
		}
		if strings.Contains(repo, "/") {
			return fmt.Errorf("GitHub 规则配置错误: 只填写仓库名(不含所有者)，收到 %q", repo)
		}
	}

	return nil
}

func (v *GithubMakeRepositoryPublicRuleValidator) Meta() meta.RuleMeta {
	return meta.RuleMeta{
		Type:        RuleTypeMakeRepoPublic,
		Label:       "公开 GitHub 仓库",
		Description: "触发时把指定的私有仓库改为公开。",
		AccountType: AccountType,
		Fields: []meta.Field{
			{
				Key:         "repos",
				Label:       "仓库列表",
				Type:        meta.FieldStringList,
				Required:    true,
				Placeholder: "my-repo",
				Help:        "只填仓库名，所有者取自关联的 GitHub 账号",
			},
		},
	}
}
