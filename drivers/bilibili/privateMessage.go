package bilibili

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CuteReimu/bilibili/v2"
	"github.com/ssdomei232/goodBaby/internal/retry"
)

// privateMsgTypeText 私信消息类型：1 为文字
const privateMsgTypeText = 1

// SendPrivateMessage 用给定账号向配置中的每个 UID 发送一条 B 站私信
//
// 每个接收者单独重试、互不阻塞，最后汇总失败信息交给上层记录。
func SendPrivateMessage(ctx context.Context, account *BiliAccount, config *BiliPrivateMessageConfig) error {
	client := newClient(account.RawCookies)

	// 发私信需要带上自己的 UID，从账号信息里取
	selfUID, err := fetchSelfUID(ctx, client)
	if err != nil {
		return err
	}

	var fails []string
	for _, uid := range config.ReceiverUids {
		if err := retry.Do(ctx, func() error {
			return sendPrivateMessage(client, selfUID, uid, config.Msg)
		}); err != nil {
			fails = append(fails, fmt.Sprintf("UID %d: %v", uid, err))
		}
	}

	if len(fails) > 0 {
		return fmt.Errorf("%d/%d 条私信发送失败: %s",
			len(fails), len(config.ReceiverUids), strings.Join(fails, "; "))
	}
	return nil
}

// fetchSelfUID 取当前账号的 UID，私信接口需要它
func fetchSelfUID(ctx context.Context, client *bilibili.Client) (int, error) {
	var selfUID int
	if err := retry.Do(ctx, func() error {
		info, err := client.GetAccountInformation()
		if err != nil {
			return err
		}
		if info == nil || info.Mid == 0 {
			return retry.Permanent(fmt.Errorf("B站 cookie 无效或已过期"))
		}
		selfUID = info.Mid
		return nil
	}); err != nil {
		return 0, fmt.Errorf("获取B站账号信息失败: %w", err)
	}
	return selfUID, nil
}

// sendPrivateMessage 发送一条私信
func sendPrivateMessage(client *bilibili.Client, senderUID int, receiverUID int64, msg string) error {
	_, err := client.SendPrivateMessage(bilibili.SendPrivateMessageParam{
		SenderUid:    senderUID,
		ReceiverId:   int(receiverUID),
		ReceiverType: 1,
		MsgType:      privateMsgTypeText,
		Timestamp:    int(time.Now().Unix()),
		Content:      buildTextContent(msg),
	})
	return err
}
