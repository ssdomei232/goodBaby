# 发版流程

发版由 [`.github/workflows/release.yml`](../.github/workflows/release.yml) 全自动完成，
触发条件是推送一个 `v` 开头的 tag。

## 发一个版本

```bash
git tag v2.2.0
git push origin v2.2.0
```

推上去之后 workflow 会依次做四件事：

| 阶段 | 内容 |
| --- | --- |
| `frontend` | 跑一次 `npm ci && npm run build`（含 `vue-tsc` 类型检查），产物作为 artifact 共享 |
| `build` | 5 个平台交叉编译，打包成 tar.gz / zip |
| `release` | 汇总产物、生成 `checksums.txt`、创建 GitHub Release（发布说明自动生成） |
| `docker` | buildx 构建 `linux/amd64` + `linux/arm64` 镜像推送到 ghcr |

前端只构建一次，所有平台的二进制共用同一份产物，保证各平台内嵌的前端完全一致。

## 产物

二进制覆盖：

* `linux/amd64`、`linux/arm64`
* `windows/amd64`
* `darwin/amd64`、`darwin/arm64`

每个压缩包内含可执行文件、`README.md`、`LICENSE` 与 `docs/`。

镜像标签由 tag 推导：`v2.2.0` 会产出 `2.2.0`、`2.2`、`2` 与 `latest`。

## 版本信息注入

编译时通过 ldflags 注入，启动日志里能看到：

```
goodBaby v2.2.0 (a1b2c3d, 2026-07-27T07:00:00Z) 启动中...
```

对应 `main.go` 里的 `version` / `gitCommit` / `buildDate` 三个变量。

## 关于交叉编译

项目依赖里没有 CGO——sqlite 用的是纯 Go 的 `glebarez/sqlite`——所以所有平台都能用
`CGO_ENABLED=0` 直接交叉编译，不需要各平台的 runner。

Dockerfile 里前端与后端构建阶段都标了 `--platform=$BUILDPLATFORM`，配合 `TARGETARCH`
做 Go 交叉编译，因此构建 arm64 镜像时不会把 npm 和 go build 丢进 QEMU 里跑。

## 手动验证

workflow 支持 `workflow_dispatch` 手动触发。手动跑时**不会创建 Release**（`release` job 有
`if: github.ref_type == 'tag'` 保护），只会构建产物并推一个 `sha-xxxxxxx` 标签的镜像，
适合在正式打 tag 前先验证流程。

## 权限

workflow 只用了内置的 `GITHUB_TOKEN`，不需要额外配置 secret：

```yaml
permissions:
  contents: write # 创建 Release
  packages: write # 推送镜像到 ghcr
```

首次推送后，包默认是私有的。如果要公开，去仓库的 Packages 页面把
`goodbaby` 这个包的可见性改成 Public。
