package email

type EmailAccountConfig struct {
	SMTPServer      string `json:"smtp_server"`
	Port            int    `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	TestDestination string `json:"test_destination"` // 用于测试邮件发送的地址
}

type EmailRule struct {
	Msg          string   `json:"msg"`
	Title        string   `json:"title"`
	Destinations []string `json:"destinations"` // 邮件地址列表
}
