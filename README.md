# goodBaby

摇篮系统

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
go build
./goodBaby
```

之后在日志中找到 Bilibili 的登陆二维码,打开Bilibili APP 扫码登陆

> [!NOTE]  
> 只有在登陆 Bilibili 后程序才会正常运行
