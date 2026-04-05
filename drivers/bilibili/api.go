package bilibili

import (
	"log"
	"os"

	"github.com/CuteReimu/bilibili/v2"
)

// Qrcode login
func HandleLoginWithQRCode(biliClient *bilibili.Client) {
	qrCode, err := biliClient.GetQRCode()
	if err != nil {
		log.Printf("获取二维码失败: %v", err)
		return
	}
	buf, _ := qrCode.Encode()
	os.WriteFile("tmp/qrcode.png", buf, 0644)
	qrCode.Print()

	result, err := biliClient.LoginWithQRCode(bilibili.LoginWithQRCodeParam{
		QrcodeKey: qrCode.QrcodeKey,
	})
	if err != nil || result.Code != 0 {
		log.Printf("登录失败: %v", err)
		return
	}

	log.Println("登录成功")
}
