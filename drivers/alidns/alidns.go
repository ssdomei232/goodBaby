package alidns

import (
	"context"

	"github.com/ssdomei232/goodBaby/model"
)

func deleteAliDNSRecord(ctx context.Context, rule *model.Rule) {
	aliDNSClient, err := GetAliDNSClient(rule)
}
