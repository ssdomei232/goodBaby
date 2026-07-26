package bilibili

const (
	// AccountType B 站账号类型标识
	AccountType = "bilibili"
	// RuleTypeDynamic 发送 B 站动态的规则类型标识
	RuleTypeDynamic = "bilibili-dynamic"
)

type BiliAccount struct {
	RawCookies string `json:"raw_cookies"`
}

type BiliDynamicConfig struct {
	Msg string `json:"msg"`
}
