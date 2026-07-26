package email

const (
	// AccountType 邮箱账号类型标识
	AccountType = "email"
	// RuleType 发送邮件的规则类型标识
	RuleType = "email"
)

// SMTP 加密方式
const (
	SecuritySSL      = "ssl"      // 隐式 TLS，一般是 465 端口
	SecuritySTARTTLS = "starttls" // 显式 TLS，一般是 587 端口
	SecurityNone     = "none"     // 不加密，一般是 25 端口
)

type EmailAccountConfig struct {
	SMTPServer string `json:"smtp_server"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	// 加密方式，留空时按 ssl 处理(兼容旧配置)
	Security string `json:"security"`
	// 发件人地址，留空时使用 Username
	From string `json:"from"`
	// 用于测试邮件发送的地址
	TestDestination string `json:"test_destination"`
}

// SecurityOrDefault 返回加密方式，兼容没有该字段的旧配置
func (c *EmailAccountConfig) SecurityOrDefault() string {
	if c.Security == "" {
		return SecuritySSL
	}
	return c.Security
}

// FromOrDefault 返回发件人地址
func (c *EmailAccountConfig) FromOrDefault() string {
	if c.From == "" {
		return c.Username
	}
	return c.From
}

type EmailRule struct {
	Msg          string   `json:"msg"`
	Title        string   `json:"title"`
	Destinations []string `json:"destinations"` // 邮件地址列表
}
