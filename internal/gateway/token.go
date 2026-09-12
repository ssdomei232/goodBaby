package gateway

import (
	"crypto/rand"
	"encoding/hex"
)

// NewToken 生成一个消息网关的访问 Token
func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "gw_" + hex.EncodeToString(b), nil
}
