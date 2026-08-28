package alidns

import (
	alidns "github.com/alibabacloud-go/alidns-20150109/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
)

func GetAliDNSClient(rule *model.Rule) (client *alidns.Client, err error) {
	var accountConfig AliDNSAccount

	// get config
	if err := db.LoadAccountConfig(rule.AccountID, &accountConfig); err != nil {
		return nil, err
	}

	// init aliyun account config
	credentialsConfig := new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(accountConfig.AK).
		SetAccessKeySecret(accountConfig.SK)
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
