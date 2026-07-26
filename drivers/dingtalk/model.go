package dingtalk

// RuleType 钉钉机器人规则类型标识
const RuleType = "dingtalk"

// DingTalkConfig 钉钉自定义机器人凭据，用户级提醒配置也用这个结构
type DingTalkConfig struct {
	AccessToken string `json:"access_token"`
	Secret      string `json:"secret"`
}

// DingTalkRuleConfig 钉钉规则配置：凭据 + 要发送的内容
type DingTalkRuleConfig struct {
	DingTalkConfig
	Title string `json:"title"`
	Msg   string `json:"msg"`
}
