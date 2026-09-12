package rainyun

import (
	"encoding/json"
	"fmt"

	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/public"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/rcs"
	"github.com/ssdomei232/rainyun-go-sdk/v2/rainyun/rgs"
)

// ParseAccountConfig 解析雨云账号配置
func ParseAccountConfig(config string) (*RainyunAccount, error) {
	if config == "" {
		return nil, fmt.Errorf("解析雨云账号配置失败: 规则没有关联账号")
	}

	var account RainyunAccount
	if err := json.Unmarshal([]byte(config), &account); err != nil {
		return nil, fmt.Errorf("解析雨云账号配置失败: %v", err)
	}
	return &account, nil
}

// ParseWorkOrderConfig 解析雨云工单规则配置
func ParseWorkOrderConfig(configJSON string) (*RainyunWorkOrderRule, error) {
	var config RainyunWorkOrderRule
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析雨云工单规则配置失败: %v", err)
	}
	return &config, nil
}

// ParseRunAwayConfig 解析雨云跑路规则配置
func ParseRunAwayConfig(configJSON string) (*RainyunRunAwayRule, error) {
	var config RainyunRunAwayRule
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("解析雨云跑路规则配置失败: %v", err)
	}
	return &config, nil
}

// reinstallOneRCS 重装一台 RCS
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
			_, err := rcsClient.ReinstallRcs(rcsID, &rcs.ReinstallRcsRequest{
				OsID:     os.ID,
				ResetOsd: true,
			})
			return err
		}
	}
	return fmt.Errorf("没有找到可用的 linux 镜像")
}

// reinstallOneRGS 重装一台 RGS
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
			_, err := rgsClient.Reinstallgs(rgsID, &rcs.ReinstallRcsRequest{
				OsID:     os.ID,
				ResetOsd: true,
			})
			return err
		}
	}
	return fmt.Errorf("没有找到可用的 linux 镜像")
}
