package onebot

const (
	// AccountType OneBot 账号类型标识
	AccountType = "onebot"
	// RuleType 发送 OneBot 消息的规则类型标识
	RuleType = "onebot"
)

type OneBotConfig struct {
	SendGroups []int64 `json:"send_groups"` // 群号列表
	SendUsers  []int64 `json:"send_users"`  // 好友 QQ 号列表
	Msg        string  `json:"msg"`
}

type OneBotAccount struct {
	URL   string `json:"url"`   // 例如 "http://localhost:5700"
	Token string `json:"token"` // 例如 "your_token_here"
}

// apiResponse OneBot HTTP API 的通用响应信封
type apiResponse struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Wording string `json:"wording"`
}

// err 把 OneBot 返回的业务错误转换成 Go error
func (r *apiResponse) err() error {
	if r.Status == "ok" && r.Retcode == 0 {
		return nil
	}
	detail := r.Wording
	if detail == "" {
		detail = r.Message
	}
	if detail == "" {
		detail = "请检查 OneBot 地址与 Token"
	}
	return &APIError{Status: r.Status, Retcode: r.Retcode, Detail: detail}
}

// APIError OneBot 返回的业务错误
type APIError struct {
	Status  string
	Retcode int
	Detail  string
}

func (e *APIError) Error() string {
	return "OneBot 返回错误(status=" + e.Status + "): " + e.Detail
}
