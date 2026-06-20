package github

import (
	"encoding/json"
	"fmt"
)

// GitHubAccountConfigValidator GitHub账号配置验证器
type GitHubAccountConfigValidator struct{}

func (v *GitHubAccountConfigValidator) GetType() string {
	return "github-account"
}

func (v *GitHubAccountConfigValidator) Validate(config string) error {
	// GitHub账号配置验证逻辑
	var cfg GithubAccount
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return fmt.Errorf("解析GitHub账号配置失败: %v", err)
	}

	if cfg.Token == "" {
		return fmt.Errorf("GitHub账号配置中Token不能为空")
	}
	if cfg.Owner == "" {
		return fmt.Errorf("GitHub账号配置中Owner不能为空")
	}

	return nil
}
