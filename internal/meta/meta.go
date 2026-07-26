// Package meta 定义驱动向前端暴露的元数据。
//
// WebUI 依靠这些描述动态渲染账号 / 规则的配置表单，
// 新增一个驱动时只要提供 Meta()，前端无需改动。
package meta

// FieldType 决定前端使用哪种输入控件
type FieldType string

const (
	FieldString     FieldType = "string"      // 单行文本
	FieldPassword   FieldType = "password"    // 密码框
	FieldTextarea   FieldType = "textarea"    // 多行文本
	FieldNumber     FieldType = "number"      // 数字
	FieldStringList FieldType = "string-list" // 字符串数组
	FieldNumberList FieldType = "number-list" // 数字数组
	FieldBool       FieldType = "bool"        // 开关
)

// Field 描述配置 JSON 中的一个字段
type Field struct {
	Key         string    `json:"key"`
	Label       string    `json:"label"`
	Type        FieldType `json:"type"`
	Required    bool      `json:"required"`
	Placeholder string    `json:"placeholder,omitempty"`
	Help        string    `json:"help,omitempty"`
	// Secret 为 true 的字段在读取接口中会被掩码
	Secret  bool `json:"secret,omitempty"`
	Default any  `json:"default,omitempty"`
}

// AccountMeta 账号类型的元数据
type AccountMeta struct {
	Type        string  `json:"type"`
	Label       string  `json:"label"`
	Description string  `json:"description,omitempty"`
	Docs        string  `json:"docs,omitempty"`
	Testable    bool    `json:"testable"`
	Fields      []Field `json:"fields"`
}

// RuleMeta 规则类型的元数据
type RuleMeta struct {
	Type        string `json:"type"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	Docs        string `json:"docs,omitempty"`
	// AccountType 该规则需要的账号类型，为空表示不需要关联账号
	AccountType string  `json:"account_type"`
	Fields      []Field `json:"fields"`
}

// SecretKeys 返回需要掩码的字段名
func SecretKeys(fields []Field) []string {
	keys := make([]string, 0, len(fields))
	for _, f := range fields {
		if f.Secret {
			keys = append(keys, f.Key)
		}
	}
	return keys
}
