# Email 配置

账号类型: `email`，规则类型: `email`。

## 账号配置 (Account Config)

```json
{
    "smtp_server": "smtp.example.com",
    "port": 465,
    "security": "ssl",
    "username": "you@example.com",
    "password": "你的密码或授权码",
    "from": "you@example.com",
    "test_destination": "test@example.com"
}
```

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `smtp_server` | 是 | SMTP 服务器地址 |
| `port` | 是 | 端口，SSL 一般为 465，STARTTLS 一般为 587 |
| `security` | 否 | `ssl` / `starttls` / `none`，默认 `ssl` |
| `username` | 是 | 登录用户名，通常是邮箱地址 |
| `password` | 是 | 密码或授权码(QQ 邮箱、163 等需要使用授权码) |
| `from` | 否 | 发件人地址，默认使用 `username` |
| `test_destination` | 否 | 填写后，点击“测试”会真的发送一封测试邮件到该地址 |

## 规则配置 (Rule Config)

```json
{
    "title": "邮件标题",
    "msg": "邮件正文",
    "destinations": ["a@example.com", "b@example.com"]
}
```

每个收件人独立发送、独立重试，互不影响。
