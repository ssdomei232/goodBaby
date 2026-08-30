package alidns

import (
	"context"

	"github.com/alibabacloud-go/alidns-20150109/v5/client"
	"github.com/ssdomei232/goodBaby/model"
)

func deleteAliDNSRecord(ctx context.Context, rule *model.Rule) error {
	aliDNSClient, err := getAliDNSClient(rule)
	if err != nil {
		return err
	}
	deleteRecordConfig, err := getDeleteRecordConfig(rule)
	if err != nil {
		return err
	}

	_, err = aliDNSClient.DeleteDomainRecord(&client.DeleteDomainRecordRequest{
		RecordId: &deleteRecordConfig.RecordID,
	})
	if err != nil {
		return err
	}
	return nil
}
