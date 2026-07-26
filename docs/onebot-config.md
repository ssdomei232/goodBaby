# OneBot (QQ) 配置

账号类型: `onebot`，规则类型: `onebot`。

对接 [NapCat](https://napneko.github.io/)、go-cqhttp 等实现了 OneBot HTTP API 的服务。

## 账号配置 (Account Config)

```json
{
    "url": "http://localhost:3000",
    "token": "your_access_token"
}
```

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `url` | 是 | OneBot HTTP 服务地址 |
| `token` | 是 | Access Token |

在 WebUI 的账号页面可以点击“测试”验证连通性(调用 `get_login_info`)。

## 规则配置 (Rule Config)

```json
{
    "msg": "要发送的消息",
    "send_groups": [123456789],
    "send_users": [987654321]
}
```

`send_groups`(群号)与 `send_users`(好友 QQ 号)至少填写一项。
