package alidns

import (
	"encoding/json"
	"fmt"

	"github.com/alibabacloud-go/alidns-20150109/v5/client"
	"github.com/ssdomei232/goodBaby/internal/meta"
)

type AliDNSAccountConfigValidator struct{}

func (v *AliDNSAccountConfigValidator) GetType() string { return AccountType }

func (v *AliDNSAccountConfigValidator) Validate(config string) error {
	cfg, err := parseAccount(config)
	if err != nil {
		return err
	}
	if cfg.AK == "" || cfg.SK == "" {
		return fmt.Errorf("阿里云账号配置不完整")
	}
	return nil
}

func parseAccount(config string) (*AliDNSAccount, error) {
	var cfg AliDNSAccount
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return nil, fmt.Errorf("解析阿里云账号配置失败: %v", err)
	}
	return &cfg, nil
}

func (v *AliDNSAccountConfigValidator) Test(config string) error {
	cfg, err := parseAccount(config)
	if err != nil {
		return err
	}

	aliDNSClient, err := getAliDNSClient(cfg.AK, cfg.SK)
	if err != nil {
		return err
	}

	_, err = aliDNSClient.DescribeDnsProductInstances(&client.DescribeDnsProductInstancesRequest{})
	if err != nil {
		return err
	}

	return nil
}

func (v *AliDNSAccountConfigValidator) Meta() meta.AccountMeta {
	return meta.AccountMeta{
		Type:        AccountType,
		Label:       "阿里云",
		Description: "阿里云国内站",
		Fields: []meta.Field{
			{Key: "ak", Label: "Access Key", Type: meta.FieldPassword, Required: true, Secret: true, Placeholder: "xxxxxxxxx"},
			{Key: "sk", Label: "Secret Key", Type: meta.FieldPassword, Required: true, Secret: true, Placeholder: "xxxxxxxxx"},
		},
	}
}
