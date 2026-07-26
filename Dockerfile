# ======================
# 前端构建阶段 (Frontend Stage)
# ======================
# 固定在构建机架构上运行：前端产物与目标架构无关，
# 否则交叉构建 arm64 镜像时 npm 要跑在 QEMU 里，慢得离谱
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend

WORKDIR /app/web/frontend

COPY web/frontend/package.json web/frontend/package-lock.json ./
RUN npm ci

COPY web/frontend/ ./
RUN npm run build

# ======================
# 后端构建阶段 (Builder Stage)
# ======================
# 同样固定在构建机架构上，用 Go 的交叉编译产出目标架构的二进制。
# 依赖里没有 CGO(sqlite 用的是纯 Go 实现)，所以交叉编译是安全的。
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
# 用前端构建产物覆盖占位的 dist
COPY --from=frontend /app/web/dist /app/web/dist

ARG VERSION="dev"
ARG BUILD_DATE="unknown"
ARG GIT_COMMIT="unknown"
# TARGETARCH 由 buildx 自动注入(amd64 / arm64 / ...)
ARG TARGETARCH
ARG TARGETOS=linux
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
    -trimpath \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.gitCommit=${GIT_COMMIT}" \
    -o /app/main .

# ======================
# 运行阶段 (Runtime Stage)
# ======================
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /app/main /app/main

# 数据(config.json / data.db)统一放在 /app/data，挂载这个目录即可
ENV GOODBABY_CONFIG=/app/data/config.json \
    GOODBABY_DB_PATH=/app/data/data.db
VOLUME ["/app/data"]

EXPOSE 8088

CMD ["/app/main"]
