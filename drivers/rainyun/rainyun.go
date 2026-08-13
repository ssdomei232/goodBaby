package rainyun

import (
	"context"
	"fmt"
	"strings"

	"github.com/ssdomei232/goodBaby/internal/retry"
	"github.com/ssdomei232/goodBaby/model"
	rain "github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/common"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/rcs"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/rgs"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/workorder"
)

func SendRainyunWorkorderMsg(ctx context.Context, rule *model.Rule) error {
	rainyunConfig, err := getRainyunConfigFromRule(rule)
	if err != nil {
		return err
	}

	rainyunAccount, err := getRainyunAccountFromRule(rule)
	if err != nil {
		return err
	}

	err = retry.Do(ctx, func() error {
		return sendRainyunWorkorderMsg(rainyunConfig, rainyunAccount)
	})

	if err != nil {
		return err
	}

	return nil
}

func sendRainyunWorkorderMsg(rainyunRuleConfig *RainyunWorkOrderRule, rainyunAccount *RainyunAccount) error {
	client := rain.NewClient(rainyunAccount.APIKey)
	workorderClient := workorder.Client{Client: client}

	_, err := workorderClient.CreateWorkOrder(&workorder.CreateWorkerorderRequest{
		Content:  rainyunRuleConfig.Msg,
		IsAuthed: false,
		IsUrgent: 1,
		Title:    rainyunRuleConfig.Title,
		Type:     "reward",
	})

	if err != nil {
		return err
	}

	return nil
}

func RainyunRunAway(ctx context.Context, rule *model.Rule) error {
	rainyunAccount, err := getRainyunAccountFromRule(rule)
	if err != nil {
		return err
	}
	client := rain.NewClient(rainyunAccount.APIKey)

	var fails []string
	total := 0

	// 重装所有RCS
	rcsClient := rcs.Client{Client: client}
	options := rain.EncodingStandardQueryParameters(1, 20)
	rcsList, err := rcsClient.GetRcsList(options)
	if err != nil {
		return err
	}
	if rcsList.Data.TotalRecords == 0 {
		return nil
	} else if rcsList.Data.TotalRecords > 20 {
		rcsTotal := rcsList.Data.TotalRecords
		options := rain.EncodingStandardQueryParameters(1, rcsTotal)
		rcsList, err = rcsClient.GetRcsList(options)
		if err != nil {
			return err
		}
	}
	for _, rcs := range rcsList.Data.Records {
		total++
		err := retry.Do(context.Background(), func() error {
			return reinstallOneRCS(rcs.ID, &rcsClient)
		})
		if err != nil {
			fails = append(fails, fmt.Sprintf("RCS %d: %v", rcs.ID, err))
		}
	}

	// 重装所有RGS
	rgsClient := rgs.Client{Client: client}
	rgsOptions := rain.EncodingStandardQueryParameters(1, 20)
	rgsList, err := rgsClient.GetRgsList(rgsOptions)
	if err != nil {
		return err
	}
	if rgsList.Data.TotalRecords == 0 {
		return nil
	} else if rgsList.Data.TotalRecords > 20 {
		rgsTotal := rgsList.Data.TotalRecords
		rgsOptions := rain.EncodingStandardQueryParameters(1, rgsTotal)
		rgsList, err = rgsClient.GetRgsList(rgsOptions)
		if err != nil {
			return err
		}
	}
	for _, rgs := range rgsList.Data.Records {
		total++
		err := retry.Do(context.Background(), func() error {
			return reinstallOneRGS(rgs.ID, &rgsClient)
		})
		if err != nil {
			fails = append(fails, fmt.Sprintf("RGS %d: %v", rgs.ID, err))
		}
	}

	if len(fails) > 0 {
		return fmt.Errorf("%d/%d 台云服务器重装失败: %s", len(fails), total, strings.Join(fails, "; "))
	}

	return nil
}
