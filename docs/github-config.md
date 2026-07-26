# GitHub 配置

账号类型: `github`，规则类型: `github-repo-public`。

定时器到期后，把指定的私有仓库设置为公开。

## 账号配置 (Account Config)

```json
{
    "token": "ghp_xxxxxxxx",
    "owner": "your-github-name"
}
```

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| `token` | 是 | Personal Access Token，需要 `repo` 权限 |
| `owner` | 是 | 仓库所有者(用户名或组织名) |

Token 在 GitHub → Settings → Developer settings → Personal access tokens 创建。

在 WebUI 的账号页面可以点击“测试”验证 Token 是否有效。

## 规则配置 (Rule Config)

```json
{
    "repos": ["repo-a", "repo-b"]
}
```

只填仓库名，不含所有者。
