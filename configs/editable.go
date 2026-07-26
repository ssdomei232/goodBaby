package configs

import "fmt"

// EditableConfig 是管理员可以在 WebUI 里修改的配置子集。
//
// 只收录改完立即生效的项。监听地址、数据库、session 密钥这类需要重启
// 或改错会导致服务起不来的配置，只能改配置文件。
type EditableConfig struct {
	EnableRegistry       bool `json:"enable_registry"`
	TimeoutDurationHours int  `json:"timeout_duration_hours"`
	CheckIntervalMinutes int  `json:"check_interval_minutes"`
	LogRetainCount       int  `json:"log_retain_count"`
}

// Editable 从完整配置中取出可编辑的部分
func (c Config) Editable() EditableConfig {
	return EditableConfig{
		EnableRegistry:       c.EnableRegistry,
		TimeoutDurationHours: c.TimeoutDurationHours,
		CheckIntervalMinutes: c.CheckIntervalMinutes,
		LogRetainCount:       c.LogRetainCount,
	}
}

// Validate 校验管理员提交的配置
func (e EditableConfig) Validate() error {
	if e.TimeoutDurationHours < 1 || e.TimeoutDurationHours > 72 {
		return fmt.Errorf("规则重试时长必须在 1~72 小时之间")
	}
	if e.CheckIntervalMinutes < 1 || e.CheckIntervalMinutes > 1440 {
		return fmt.Errorf("检查间隔必须在 1~1440 分钟之间")
	}
	if e.LogRetainCount < 0 || e.LogRetainCount > 100000 {
		return fmt.Errorf("日志保留条数必须在 0~100000 之间，0 表示不限制")
	}
	return nil
}

// ApplyEditable 把可编辑部分合并进当前配置并持久化，返回更新后的完整配置
func ApplyEditable(e EditableConfig) (Config, error) {
	if err := e.Validate(); err != nil {
		return Config{}, err
	}

	current, err := GetConfig()
	if err != nil {
		return Config{}, err
	}

	current.EnableRegistry = e.EnableRegistry
	current.TimeoutDurationHours = e.TimeoutDurationHours
	current.CheckIntervalMinutes = e.CheckIntervalMinutes
	current.LogRetainCount = e.LogRetainCount

	if err := Save(current); err != nil {
		return Config{}, err
	}
	return current, nil
}
