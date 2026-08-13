package rainyun

import (
	"context"

	"github.com/ssdomei232/goodBaby/model"
	rain "github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/common"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/workorder"
)

func sendRainyunWorkorderMsg(ctx context.Context, rule *model.Rule) error {
	rainyunConfig, err := getRainyunConfigFromRule(rule)
	if err != nil {
		return err
	}

	rainyunAccount, err := getRainyunAccountFromRule(rule)
	if err != nil {
		return err
	}

	client := rain.NewClient(rainyunAccount.APIKey)
	workorderClient := workorder.Client{Client: client}

	_, err = workorderClient.CreateWorkOrder(&workorder.CreateWorkerorderRequest{
		Content:  rainyunConfig.Msg,
		IsAuthed: false,
		IsUrgent: 1,
		Title:    rainyunConfig.Title,
		Type:     "reward",
	})

	if err != nil {
		return err
	}

	return nil
}
