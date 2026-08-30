package alidns

import (
	"encoding/json"

	alidns "github.com/alibabacloud-go/alidns-20150109/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

// getAliDNSClient initializes and returns an Alibaba Cloud DNS client based on the provided rule's account configuration. It retrieves the account configuration from the database, sets up the necessary credentials, and creates a new client instance for interacting with Alibaba Cloud DNS services.
func getAliDNSClientFromRule(rule *model.Rule) (client *alidns.Client, err error) {
	var accountConfig AliDNSAccount

	// get config
	if err := db.LoadAccountConfig(rule.AccountID, &accountConfig); err != nil {
		return nil, err
	}

	return getAliDNSClient(accountConfig.AK, accountConfig.SK)
}

func getAliDNSClient(ak string, sk string) (client *alidns.Client, err error) {
	// init aliyun account config
	credentialsConfig := new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(ak).
		SetAccessKeySecret(sk)
	akCredential, err := credentials.NewCredential(credentialsConfig)
	if err != nil {
		return nil, err
	}
	config := &openapi.Config{}
	config.Credential = akCredential
	config.Endpoint = tea.String("alidns.aliyuncs.com")

	// create aliyun account client
	client, _ = alidns.NewClient(config)

	return client, nil
}

// getDeleteRecordConfig retrieves the configuration for deleting DNS records from the provided rule. It unmarshals the rule's configuration JSON into a DeleteRecordConfig structure and returns it.
func getDeleteRecordConfig(rule *model.Rule) (*DeleteRecordConfig, error) {
	var deleteRecordConfig DeleteRecordConfig

	if err := json.Unmarshal([]byte(rule.ConfigJson), &deleteRecordConfig); err != nil {
		return nil, err
	}

	return &deleteRecordConfig, nil
}
