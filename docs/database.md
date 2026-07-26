# 数据库配置

goodBaby 默认使用 SQLite，零配置开箱即用；也可以切换到 PostgreSQL。

## SQLite（默认）

```json
{
  "database_driver": "sqlite",
  "database_path": "data.db"
}
```

用的是纯 Go 的 [glebarez/sqlite](https://github.com/glebarez/sqlite)，不需要 CGO，编译出来就是一个静态二进制。

已开启 WAL 与 10 秒 busy timeout；因为 SQLite 是单写模型，连接池被限制为 1 条连接，避免并发执行规则时出现 `database is locked`。

## PostgreSQL

```json
{
  "database_driver": "postgres",
  "database_dsn": "postgres://user:password@localhost:5432/goodbaby?sslmode=disable"
}
```

切换到 postgres 后 `database_path` 会被忽略。连接池为最大 20 条连接、5 条空闲、连接最长存活 1 小时。

启动时会先 `Ping` 一次，配置写错会在启动阶段直接报错退出，而不是等到第一次查询才失败。

### 环境变量

配置文件之外也可以用环境变量覆盖，适合容器部署：

| 变量 | 说明 |
| --- | --- |
| `GOODBABY_DB_DRIVER` | `sqlite` 或 `postgres` |
| `GOODBABY_DB_DSN` | postgres 连接串 |
| `GOODBABY_DB_PATH` | sqlite 文件路径 |

只设置了 `GOODBABY_DB_DSN` 而没设置 `GOODBABY_DB_DRIVER` 时，会自动按 postgres 处理。

### docker compose 示例

```yaml
services:
  goodbaby:
    image: mei232/goodbaby:v2.1.0
    ports:
      - "8088:8088"
    environment:
      - TZ=Asia/Shanghai
      - GIN_MODE=release
      - GOODBABY_DB_DSN=postgres://goodbaby:goodbaby@postgres:5432/goodbaby?sslmode=disable
    volumes:
      - ./data:/app/data
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

  postgres:
    image: postgres:17-alpine
    environment:
      - POSTGRES_USER=goodbaby
      - POSTGRES_PASSWORD=goodbaby
      - POSTGRES_DB=goodbaby
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U goodbaby"]
      interval: 5s
      timeout: 5s
      retries: 10
    restart: unless-stopped

volumes:
  pgdata:
```

## 切换数据库

两种数据库之间**没有**自动数据迁移。切换驱动后是一个全新的空库，需要重新注册账号并配置规则。

如果要搬运已有数据，得自己从 SQLite 导出再导入 PostgreSQL（表结构一致：`users` / `timers` / `rules` / `accounts` / `execution_logs`）。

## 建表

首次启动会自动执行 `AutoMigrate` 建表并补齐新增字段，不需要手动执行 SQL。
