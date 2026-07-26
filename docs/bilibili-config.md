# Bilibili 配置

账号类型: `bilibili`，规则类型: `bilibili-dynamic`。

## 账号配置 (Account Config)

```json
{
    "raw_cookies": "SESSDATA=xxx; bili_jct=xxx; DedeUserID=xxx; ..."
}
```

登录 [bilibili.com](https://www.bilibili.com) 后，从浏览器开发者工具(F12 → 网络 → 任意请求 → 请求头 → Cookie)复制完整的 Cookie 字符串。

在 WebUI 的账号页面可以点击“测试”验证 Cookie 是否有效。

## 规则配置 (Rule Config)

```json
{
    "msg": "要发送的动态内容"
}
```

定时器到期后会以该账号发布一条纯文本动态。
