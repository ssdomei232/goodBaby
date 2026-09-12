package alidns

import (
	"encoding/json"
	"fmt"

	alidns "github.com/alibabacloud-go/alidns-20150109/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
)

// ParseAccountConfig 解析阿里云账号配置
func ParseAccountConfig(config string) (*AliDNSAccount, error) {
	if config == "" {
		return nil, fmt.Errorf("解析阿里云账号配置失败: 规则没有关联账号")
	}

	var account AliDNSAccount
	if err := json.Unmarshal([]byte(config), &account); err != nil {
		return nil, fmt.Errorf("解析阿里云账号配置失败: %v", err)
	}
	return &account, nil
}

// ParseDeleteRecordConfig 解析删除解析记录规则配置
func ParseDeleteRecordConfig(configJSON string) (*DeleteRecordConfig, error) {
	var config DeleteRecordConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析阿里云删除记录规则配置失败: %v", err)
	}
	return &config, nil
}

// newClient 用 AccessKey 创建阿里云 DNS 客户端
func newClient(ak string, sk string) (*alidns.Client, error) {
	credentialsConfig := new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(ak).
		SetAccessKeySecret(sk)

	akCredential, err := credentials.NewCredential(credentialsConfig)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云凭据失败: %w", err)
	}

	return alidns.NewClient(&openapi.Config{
		Credential: akCredential,
		Endpoint:   tea.String("alidns.aliyuncs.com"),
	})
}
