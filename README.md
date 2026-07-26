# goodBaby v2

摇篮系统 —— 为独居青年准备的「死人开关」(Dead Man's Switch)。

一个人在大城市生活，最怕的是出了事都没人知道。goodBaby 让你定期回来报一声平安；一旦超过设定的期限没有签到，它就会替你把预先托付的留言送出去：给家人的邮件、给朋友的 QQ / 钉钉消息、一条 B 站动态，或是公开你的 GitHub 仓库，把作品留给世界。

[![通过雨云一键部署](https://rainyun-apps.cn-nb1.rains3.com/materials/deploy-on-rainyun-cn.svg)](https://app.rainyun.com/apps/rca/store/7125/cat_)

## 功能

* **WebUI**：内置 Vue3 前端，注册 / 登录、定时器、规则、账号、执行日志、设置全部可视化操作
* **定时器 (Timer)**：设定签到周期与提前提醒时间，到期前通过钉钉机器人提醒你签到
* **规则 (Rule)**：定时器到期后要执行的动作，支持：
  * 发送邮件 (SMTP)
  * 发布 B 站动态
  * 发送 QQ 消息 (OneBot / NapCat)
  * 发送钉钉机器人消息
  * 公开 GitHub 仓库
* **账号 (Account)**：集中管理第三方凭据，支持连通性测试，敏感字段(密码/Cookie/Token)不会回显
* **执行日志**：每次规则执行与提醒都有记录，规则支持手动测试
* **配置测试**：账号可一键测试连通性；规则可手动触发验证

## 快速开始

### Docker Compose (推荐)

```bash
docker compose up -d
```

首次启动会在 `./data` 下自动生成 `config.json` 与 `data.db`。
打开 `http://localhost:8088`，按引导创建第一个账号即可。

### 从源码构建

需要 Go 1.25+ 与 Node.js 20+：

```bash
# 1. 构建前端(产物输出到 web/dist，会被 go:embed 打进二进制)
cd web/frontend && npm install && npm run build && cd ../..

# 2. 构建后端
go build -o goodbaby .

# 3. 运行
./goodbaby
```

### 本地开发

```bash
# 终端 1: 启动后端 (监听 :8088)
go run .

# 终端 2: 启动前端 dev server (监听 :5173，API 代理到 :8088)
cd web/frontend && npm run dev
```

## 配置

配置文件默认为工作目录下的 `config.json`，首次启动自动生成，字段均有默认值：

| 字段 | 默认值 | 说明 |
| --- | --- | --- |
| `listen_addr` | `:8088` | HTTP 监听地址 |
| `enable_registry` | `true` | 是否开放注册（系统无用户时始终允许注册第一个账号） |
| `timeout_duration_hours` | `6` | 规则执行失败后指数退避重试的最长时间(小时) |
| `check_interval_minutes` | `10` | 检查定时器的间隔(分钟) |
| `database_path` | `data.db` | sqlite 数据库路径 |
| `session_secret` | 自动生成 | 会话加密密钥，自动生成并持久化 |
| `session_max_age_hours` | `168` | 会话有效期(小时) |
| `allowed_origins` | `[]` | 允许跨域的来源，前端本地开发时可填 `["http://localhost:5173"]` |
| `log_retain_count` | `500` | 每个用户保留的执行日志条数 |

环境变量覆盖：`GOODBABY_CONFIG`(配置文件路径)、`GOODBABY_LISTEN_ADDR`、`GOODBABY_DB_PATH`、`GOODBABY_SESSION_SECRET`、`GOODBABY_ENABLE_REGISTRY`。

## 驱动配置

各驱动的账号 / 规则配置说明见 [docs/](docs/)：

* [Bilibili](docs/bilibili-config.md)
* [Email](docs/email-config.md)
* [OneBot (QQ)](docs/onebot-config.md)
* [钉钉机器人](docs/dingtalk-config.md)
* [GitHub](docs/github-config.md)

## 开发：新增一种规则类型

1. 在 `drivers/<name>/` 下实现：
   * 规则验证器（实现 `ruleConfigChecker.RuleValidator`，`Meta()` 返回表单元数据）
   * 执行器（实现 `runner.RuleExecutor`）
   * 如需第三方凭据，再实现账号验证器（`accountConfigChecker.AccountValidator`，可选实现 `AccountTester` 支持连通性测试）
2. 在 `internal/ruleConfigChecker/reg.go`、`internal/accountConfigChecker/reg.go`、`handler/runner/interface.go` 中注册
3. 前端无需改动 —— WebUI 会根据 `Meta()` 返回的字段描述自动渲染配置表单
