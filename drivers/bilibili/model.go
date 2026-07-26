package bilibili

const (
	// AccountType B 站账号类型标识
	AccountType = "bilibili"
	// RuleTypeDynamic 发送 B 站动态的规则类型标识
	RuleTypeDynamic = "bilibili-dynamic"
	// RuleTypePrivateMessage 发送 B 站私信的规则类型标识
	RuleTypePrivateMessage = "bilibili-private-message"
)

type BiliAccount struct {
	RawCookies string `json:"raw_cookies"`
}

type BiliDynamicConfig struct {
	Msg string `json:"msg"`
}

// BiliPrivateMessageConfig 私信规则配置
type BiliPrivateMessageConfig struct {
	// 接收者的 UID 列表
	ReceiverUids []int64 `json:"receiver_uids"`
	Msg          string  `json:"msg"`
}
