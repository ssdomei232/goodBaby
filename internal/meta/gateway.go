package meta

import "encoding/json"

// GatewayMessagePlaceholder 是「由投递内容提供」的字段占位值。
//
// 这些字段在「网关规则」页面不由用户填写，真实内容来自 webhook 请求；
// 但规则配置里必须保留该字段，投递时才知道要覆盖哪一个。
const GatewayMessagePlaceholder = "由消息网关投递内容替换"

// GatewayMessageFields 返回所有由投递内容提供的字段
func GatewayMessageFields(fields []Field) []Field {
	matched := make([]Field, 0, len(fields))
	for _, f := range fields {
		if f.GatewayMessage {
			matched = append(matched, f)
		}
	}
	return matched
}

// FillGatewayMessages 给配置里「由投递内容提供」的字段补上占位值。
//
// 网关规则提交的配置里通常没有 title / msg 这类字段(前端不让填)，
// 校验前先补齐，规则才会既通过校验、又能在投递时被覆盖。
func FillGatewayMessages(configJSON string, fields []Field) string {
	messageFields := GatewayMessageFields(fields)
	if len(messageFields) == 0 || configJSON == "" {
		return configJSON
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(configJSON), &raw); err != nil {
		return configJSON
	}

	changed := false
	for _, f := range messageFields {
		if value, ok := raw[f.Key].(string); ok && value != "" {
			continue
		}
		raw[f.Key] = GatewayMessagePlaceholder
		changed = true
	}
	if !changed {
		return configJSON
	}

	filled, err := json.Marshal(raw)
	if err != nil {
		return configJSON
	}
	return string(filled)
}
