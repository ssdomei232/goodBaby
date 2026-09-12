package alidns

import (
	"context"
	"fmt"

	alidns "github.com/alibabacloud-go/alidns-20150109/v5/client"
	"github.com/ssdomei232/goodBaby/internal/retry"
)

// DeleteRecord 用给定账号删除一条阿里云 DNS 解析记录
func DeleteRecord(ctx context.Context, account *AliDNSAccount, config *DeleteRecordConfig) error {
	client, err := newClient(account.AK, account.SK)
	if err != nil {
		return err
	}

	if err := retry.Do(ctx, func() error {
		_, err := client.DeleteDomainRecord(&alidns.DeleteDomainRecordRequest{
			RecordId: &config.RecordID,
		})
		return err
	}); err != nil {
		return fmt.Errorf("删除阿里云DNS记录失败: %w", err)
	}

	return nil
}
