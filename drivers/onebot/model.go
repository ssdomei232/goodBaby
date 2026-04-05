package onebot

type OneBotConfig struct {
	SendGroups []int  `json:"send_groups"`
	Msg        string `json:"msg"`
}

type OneBotAccount struct {
	URL   string `json:"url"`   // 例如 "http://localhost:5700"
	Token string `json:"token"` // 例如 "your_token_here"
}
