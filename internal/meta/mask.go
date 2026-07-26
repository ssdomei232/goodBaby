package meta

import "encoding/json"

// MaskPlaceholder 是掩码后的占位值。
// 提交编辑时如果某个 secret 字段仍是这个值，服务端会保留原值。
const MaskPlaceholder = "********"

// Mask 把配置 JSON 中的敏感字段替换为占位符
func Mask(configJSON string, fields []Field) string {
	secrets := SecretKeys(fields)
	if len(secrets) == 0 || configJSON == "" {
		return configJSON
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(configJSON), &raw); err != nil {
		return configJSON
	}

	for _, key := range secrets {
		if v, ok := raw[key]; ok {
			if s, isStr := v.(string); !isStr || s == "" {
				continue
			}
			raw[key] = MaskPlaceholder
		}
	}

	masked, err := json.Marshal(raw)
	if err != nil {
		return configJSON
	}
	return string(masked)
}

// Unmask 用旧配置补回被掩码的字段，返回可以落库的配置 JSON
func Unmask(newConfigJSON, oldConfigJSON string, fields []Field) string {
	secrets := SecretKeys(fields)
	if len(secrets) == 0 || newConfigJSON == "" || oldConfigJSON == "" {
		return newConfigJSON
	}

	var newRaw, oldRaw map[string]any
	if err := json.Unmarshal([]byte(newConfigJSON), &newRaw); err != nil {
		return newConfigJSON
	}
	if err := json.Unmarshal([]byte(oldConfigJSON), &oldRaw); err != nil {
		return newConfigJSON
	}

	restored := false
	for _, key := range secrets {
		if s, ok := newRaw[key].(string); ok && s == MaskPlaceholder {
			if old, exists := oldRaw[key]; exists {
				newRaw[key] = old
				restored = true
			}
		}
	}
	if !restored {
		return newConfigJSON
	}

	merged, err := json.Marshal(newRaw)
	if err != nil {
		return newConfigJSON
	}
	return string(merged)
}
