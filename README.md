# goodBaby

摇篮系统  

[![通过雨云一键部署](https://rainyun-apps.cn-nb1.rains3.com/materials/deploy-on-rainyun-cn.svg)](https://app.rainyun.com/apps/rca/store/7125/cat_)

## 他会做什么

1. 接受 singnal 请求
2. 当定时器超时(过长时间没有收到请求)时，触发 `trigger` 函数
3. 发送消息到QQ(依赖[CatBot](https://github.com/ssdomei232/CatBot))
4. 发送邮件
5. 将 Github 上的遗言仓库设置为 public
6. 发送一条 Bilibili 动态

## 配置

参照 config.example.json 修改 config.json

```json
{
    "basic": {
        "name": "xx",               // 你的名字
        "qq_number": "123",         // QQ号
        "age": 16,                  // 你的年龄
        "cause_stop": "自杀"        // 你最可能的死因
    },
    "signal_secret": "",            // 发送 signal请求时需要的密钥
    "debug": false,                 // debug 模式
    "disconnect_duration": 5,       // 失连时间，超出后会触发系统
    "enable_qq": true,              // 是否启用QQ
    "cat_bot_url": "",              // CatBot 的 webhook地址(https://github.com/ssdomei232/CatBot)
    "cat_bot_key": "",              // CatBot 的 key
    "qq_send_group": [],            // 要发送QQ消息的群号
    "qq_msg": "",                   // 要附加的QQ消息
    "mail_list": [],                // 要发送的邮箱
    "mail_title": "摇篮系统已触发",   // 邮件标题
    "smtp_config": {                // 邮件服务器配置
        "host": "",
        "port": 465,
        "user": "",
        "pass": ""
    },
    "mail_content": "",             // 要附加的邮件内容
    "github_config": {              // Github 配置
        "owner": "",
        "repos": [""],
        "token": "ghp_xxxx"
    },
    "bili_msg": "xxx",               // 哔哩哔哩动态的附加消息
    "bili_warn_address": ""          // 哔哩哔哩cookie过期发送告警的邮件地址
}
```

## 运行

```bash
mkdir goodbaby
cd goodbaby
wget https://raw.githubusercontent.com/ssdomei232/goodBaby/refs/heads/main/docker-compose.yml
docker compose up -d
docker logs goodbaby ## 扫码登陆后才能使用
```

在日志中找到 Bilibili 的登陆二维码或在文件管理中找到`tmp/qrcode.png`,打开Bilibili APP 扫码登陆

> [!NOTE]  
> 只有在登陆 Bilibili 后程序才会正常运行

## 使用

定期请求siganl接口,示例:`http://192.168.1.245:8088/signal?secret=xxxx`
