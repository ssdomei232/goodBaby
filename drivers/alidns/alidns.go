package alidns

import (
	"context"
	"fmt"

	"github.com/alibabacloud-go/alidns-20150109/v5/client"
	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
)

// deleteAliDNSRecord deletes a DNS record from Alibaba Cloud DNS based on the provided rule.
func deleteAliDNSRecord(ctx context.Context, rule *model.Rule) error {
	aliDNSClient, err := getAliDNSClientFromRule(rule)
	if err != nil {
		return err
	}
	deleteRecordConfig, err := getDeleteRecordConfig(rule)
	if err != nil {
		return err
	}

	if err := retry.Do(ctx, func() error {
		_, err := aliDNSClient.DeleteDomainRecord(&client.DeleteDomainRecordRequest{
			RecordId: &deleteRecordConfig.RecordID,
		})
		return err
	}); err != nil {
		return fmt.Errorf("删除阿里云DNS记录失败: %w", err)
	}

	return nil
}
