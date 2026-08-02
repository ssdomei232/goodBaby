# 饭碗警告配置

规则类型: `fwalert`。

定时器到期后，请求饭碗警告的 Webhook 地址

饭碗警告: 比较简单的(不用实名认证)支持电话和短信通知的第三方平台: `https://fwalert.com/`,价格不算最低，但比较方便  

## 规则配置 (Rule Config)

```json
{
    "webhook_url": "https://fwalert.com/xxx-xxx-xxx-xxx-xxx",
    "msg": "消息内容"
}
```

## 饭碗警告平台配置

goodbaby 在请求饭碗警告时将`msg`放在请求体中，参照下面和官方教程设置:  

| 字段 | 值 |
| --- | --- |
| `变量名` | `message` |
| `变量来源` | `请求体（body）` |
| `提取规则` | `序列化数据` |
| `键` | `message` |
| `通知简述` | 随便填 |
| `通知正文` | `{{message}}` |
| `联系方式类型` | 参照官方文档添加并选择 |
