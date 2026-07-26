# Bilibili 配置

账号类型: `bilibili`。

支持的规则类型:

| 规则类型 | 说明 |
| --- | --- |
| `bilibili-dynamic` | 发布一条纯文本动态 |
| `bilibili-private-message` | 向指定用户发送私信 |

## 账号配置 (Account Config)

```json
{
    "raw_cookies": "SESSDATA=xxx; bili_jct=xxx; DedeUserID=xxx; ..."
}
```

登录 [bilibili.com](https://www.bilibili.com) 后，从浏览器开发者工具(F12 → 网络 → 任意请求 → 请求头 → Cookie)复制完整的 Cookie 字符串。

在 WebUI 的账号页面可以点击“测试”验证 Cookie 是否有效。

## 动态规则配置 (`bilibili-dynamic`)

```json
{
    "msg": "要发送的动态内容"
}
```

定时器到期后会以该账号发布一条纯文本动态。

## 私信规则配置 (`bilibili-private-message`)

```json
{
    "receiver_uids": [2, 123456],
    "msg": "要发送的私信内容"
}
```

`receiver_uids` 是接收者的 UID 列表，即对方空间地址 `space.bilibili.com/` 后面的那串数字。

发送者 UID 由程序自动从账号 Cookie 中读取，不需要填写。

注意：B 站对私信有频率限制与陌生人私信限制，如果对方设置了“仅关注的人可以给我发私信”，发送会失败。失败信息可以在执行日志里看到。

## 关于视频投稿

当前使用的 SDK [`CuteReimu/bilibili`](https://github.com/CuteReimu/bilibili) **没有提供视频投稿(上传稿件)接口**，因此暂不支持“定时投稿视频”这类规则。

该 SDK 与上传相关的能力只有 `UploadDynamicBfs`（为图片动态上传图片），不包含投稿所需的 `preupload` / 分片上传 / `archive/add` 流程。

若要支持视频投稿，需要:

1. 自行实现 B 站投稿协议（预上传 → 分片上传 → 合并 → 提交稿件），且该协议未公开、变动频繁；
2. 给 goodBaby 增加文件上传与存储能力（当前只存 JSON 配置，不保存任何文件）；
3. 处理大文件的断点续传与磁盘占用。

这是一个独立的功能模块，不在现有驱动的改造范围内。
