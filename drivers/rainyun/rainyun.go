package rainyun

import (
	"context"
	"fmt"
	"strings"

	"github.com/ssdomei232/goodBaby/internal/retry"
	rain "github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/common"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/rcs"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/rgs"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/workorder"
)

// SendWorkOrder 用给定账号发送一条雨云工单
func SendWorkOrder(ctx context.Context, account *RainyunAccount, config *RainyunWorkOrderRule) error {
	if err := retry.Do(ctx, func() error {
		return createWorkOrder(account, config)
	}); err != nil {
		return fmt.Errorf("发送雨云工单失败: %w", err)
	}
	return nil
}

// createWorkOrder 调用一次工单接口
func createWorkOrder(account *RainyunAccount, config *RainyunWorkOrderRule) error {
	workorderClient := workorder.Client{Client: rain.NewClient(account.APIKey)}

	_, err := workorderClient.CreateWorkOrder(&workorder.CreateWorkerorderRequest{
		Content:  config.Msg,
		IsAuthed: false,
		IsUrgent: 1,
		Title:    config.Title,
		Type:     "reward",
	})
	return err
}

// RunAway 重装账号下所有云服务器(一键跑路)
func RunAway(ctx context.Context, account *RainyunAccount) error {
	client := rain.NewClient(account.APIKey)

	rcsFails, rcsTotal, err := reinstallAllRcs(ctx, client)
	if err != nil {
		return err
	}
	rgsFails, rgsTotal, err := reinstallAllRgs(ctx, client)
	if err != nil {
		return err
	}

	fails := append(rcsFails, rgsFails...)
	total := rcsTotal + rgsTotal
	if len(fails) > 0 {
		return fmt.Errorf("%d/%d 台云服务器重装失败: %s", len(fails), total, strings.Join(fails, "; "))
	}
	return nil
}

// reinstallAllRcs 重装全部 RCS，返回失败信息与云服务器总数
func reinstallAllRcs(ctx context.Context, client *rain.Client) ([]string, int, error) {
	rcsClient := rcs.Client{Client: client}

	ids, err := listRcsIDs(&rcsClient)
	if err != nil {
		return nil, 0, err
	}

	var fails []string
	for _, id := range ids {
		if err := retry.Do(ctx, func() error {
			return reinstallOneRCS(id, &rcsClient)
		}); err != nil {
			fails = append(fails, fmt.Sprintf("RCS %d: %v", id, err))
		}
	}
	return fails, len(ids), nil
}

// reinstallAllRgs 重装全部 RGS，返回失败信息与云服务器总数
func reinstallAllRgs(ctx context.Context, client *rain.Client) ([]string, int, error) {
	rgsClient := rgs.Client{Client: client}

	ids, err := listRgsIDs(&rgsClient)
	if err != nil {
		return nil, 0, err
	}

	var fails []string
	for _, id := range ids {
		if err := retry.Do(ctx, func() error {
			return reinstallOneRGS(id, &rgsClient)
		}); err != nil {
			fails = append(fails, fmt.Sprintf("RGS %d: %v", id, err))
		}
	}
	return fails, len(ids), nil
}

// listRcsIDs 取账号下全部 RCS 的 ID，超过一页时重新按总数拉取
func listRcsIDs(client *rcs.Client) ([]int, error) {
	first, err := client.GetRcsList(rain.EncodingStandardQueryParameters(1, pageSize))
	if err != nil {
		return nil, err
	}

	ids := make([]int, 0, len(first.Data.Records))
	for _, server := range first.Data.Records {
		ids = append(ids, server.ID)
	}
	if first.Data.TotalRecords <= pageSize {
		return ids, nil
	}

	all, err := client.GetRcsList(rain.EncodingStandardQueryParameters(1, first.Data.TotalRecords))
	if err != nil {
		return nil, err
	}

	ids = ids[:0]
	for _, server := range all.Data.Records {
		ids = append(ids, server.ID)
	}
	return ids, nil
}

// listRgsIDs 取账号下全部 RGS 的 ID，超过一页时重新按总数拉取
func listRgsIDs(client *rgs.Client) ([]int, error) {
	first, err := client.GetRgsList(rain.EncodingStandardQueryParameters(1, pageSize))
	if err != nil {
		return nil, err
	}

	ids := make([]int, 0, len(first.Data.Records))
	for _, server := range first.Data.Records {
		ids = append(ids, server.ID)
	}
	if first.Data.TotalRecords <= pageSize {
		return ids, nil
	}

	all, err := client.GetRgsList(rain.EncodingStandardQueryParameters(1, first.Data.TotalRecords))
	if err != nil {
		return nil, err
	}

	ids = ids[:0]
	for _, server := range all.Data.Records {
		ids = append(ids, server.ID)
	}
	return ids, nil
}
