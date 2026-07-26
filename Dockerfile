# ======================
# 前端构建阶段 (Frontend Stage)
# ======================
FROM node:22-alpine AS frontend

WORKDIR /app/web/frontend

COPY web/frontend/package.json web/frontend/package-lock.json ./
RUN npm ci

COPY web/frontend/ ./
RUN npm run build

# ======================
# 后端构建阶段 (Builder Stage)
# ======================
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
# 用前端构建产物覆盖占位的 dist
COPY --from=frontend /app/web/dist /app/web/dist

ARG VERSION="2.0.0"
ARG BUILD_DATE="unknown"
ARG GIT_COMMIT="unknown"
RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
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
