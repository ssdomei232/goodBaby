package rainyun

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/goodBaby/handler/db"
	"github.com/ssdomei232/goodBaby/model"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/public"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/rcs"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/rgs"
)

// 从 Rule 中获取 RainyunConfig
func getRainyunConfigFromRule(rule *model.Rule) (*RainyunWorkOrderRule, error) {
	var rainyunConfig RainyunWorkOrderRule
	if err := json.Unmarshal([]byte(rule.ConfigJson), &rainyunConfig); err != nil {
		return nil, fmt.Errorf("解析Rainyun规则配置失败: %w", err)
	}
	return &rainyunConfig, nil
}

// 从 Rule 中获取 RainyunAccount
func getRainyunAccountFromRule(rule *model.Rule) (*RainyunAccount, error) {
	var rainyunAccount RainyunAccount
	if err := db.LoadAccountConfig(rule.AccountID, &rainyunAccount); err != nil {
		return nil, err
	}
	return &rainyunAccount, nil
}

func reinstallOneRCS(rcsID int, rcsClient *rcs.Client) error {
	rcsDetail, err := rcsClient.GetRcsDetails(rcsID)
	if err != nil {
		return err
	}

	rcsOSList, err := public.GetRcsOSList()
	if err != nil {
		return err
	}

	for _, os := range rcsOSList.Data {
		if os.Region == rcsDetail.Data.Data.Node.Region && os.OsType == "linux" {
			rcsReinstallRequest := &rcs.ReinstallRcsRequest{
				OsID:     os.ID,
				ResetOsd: true,
			}
			_, err := rcsClient.ReinstallRcs(rcsID, rcsReinstallRequest)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func reinstallOneRGS(rgsID int, rgsClient *rgs.Client) error {
	rgsDetail, err := rgsClient.GetRgsDetails(rgsID)
	if err != nil {
		return err
	}

	rgsOSList, err := public.GetRgsOSList()
	if err != nil {
		return err
	}

	for _, os := range rgsOSList.Data {
		if os.Region == rgsDetail.Data.Data.Node.Region && os.OsType == "linux" {
			rgsReinstallRequest := &rcs.ReinstallRcsRequest{
				OsID:     os.ID,
				ResetOsd: true,
			}
			_, err := rgsClient.Reinstallgs(rgsID, rgsReinstallRequest)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}
