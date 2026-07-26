# 钉钉机器人配置

钉钉在 goodBaby 中有两个用途：

1. **提醒渠道**：定时器临近到期时提醒你签到，在 WebUI 的“设置”页配置
2. **规则类型** (`dingtalk`)：定时器到期后向钉钉群发送消息，不需要关联账号

## 获取机器人凭据

1. 在钉钉群中：群设置 → 智能群助手 → 添加机器人 → 自定义
2. 安全设置推荐选择“加签”，记下 Secret
3. 创建后复制 Webhook 地址中 `access_token=` 后面的值

## 规则配置 (Rule Config)

```json
{
    "access_token": "机器人 access_token",
    "secret": "加签 Secret(未开启加签可留空)",
    "title": "消息标题",
    "msg": "消息内容，支持 Markdown"
}
```

## 提醒渠道配置(设置页)

```json
{
    "access_token": "机器人 access_token",
    "secret": "加签 Secret"
}
```
