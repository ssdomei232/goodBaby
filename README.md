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
7. 当倒计时小于48小时时，向用钉钉机器人发消息提醒你

## 配置

> [!TIP]
> 更详细的配置注解请前往 [Wiki 页面](https://github.com/ssdomei232/goodBaby/wiki)

参照 config.example.json 修改 config.json

```json
{
    "signal_secret": "xxx",                                 // 发送 signal请求时需要的密钥
    "debug": false,                                         // debug 模式
    "disconnect_duration": 80,                              // 失连时间，超出后会触发系统(小时)
    "bili_msg": "摇篮系统测试",                               // 哔哩哔哩动态的附加消息
    "mail_config": {
        "host": "smtp.example.com",
        "port": 465,
        "user": "xxx",
        "pass": "xxx",
        "mail_list": ["xxx@example.com", "xxx1@example.com"],// 要发送的邮箱
        "mail_title": "摇篮系统已触发",                        // 邮件标题
        "mail_content": "这是摇篮系统的邮件内容。"              // 要附加的邮件内容
    },                                                      // 邮件配置
    "basic": {
        "name": "xx",                                       // 你的名字
        "qq_number": "123",                                 // 你的QQ号
        "nickname": "xx_123",                               // 网名
        "age": 16,                                          // 你的年龄
        "cause_stop": "自杀"                                 // 你最可能的死因
    },
    "github_config": {
        "owner": "example",                                 // github用户名
        "repos": ["example-repo1", "example-repo2"],        // 要公开的 github 仓库
        "token": "ghp_xxxx"                                 // github token
    },
    "dingtalk_bot": {
        "access_token": "xxx",
        "secret": "xxx"
    },                                                     // 钉钉机器人（用于提醒发送signal和bilibili过期提醒）
    "onebot_config": {
        "url": "http://localhost:5700",
        "token": "xxx",
        "send_groups": [12345678, 87654321],                // 要发送到的Onebot群组
        "msg": "这是OneBot的消息内容。"                       // Onebot消息
    }
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

访问你开放的端口(默认为8088),进入WebUI,输入你的密钥,点击保存(密钥会存储到浏览器的localStorage中),此时你可以在WebUI查看倒计时和发送Signal

## Todo list

- [ ] 熔断机制
- [ ] Webhook 自定义请求
- [ ] 保存定时器状态
