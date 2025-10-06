# ======================
# 构建阶段 (Builder Stage)
# ======================
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG VERSION="1.0.0"
ARG BUILD_DATE="unknown"
ARG GIT_COMMIT="unknown"
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildDate=${BUILD_DATE} -X main.gitCommit=${GIT_COMMIT}" \
    -o /app/main .

# ======================
# 运行阶段 (Runtime Stage)
# ======================
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main /app/main

EXPOSE 8088

CMD ["/app/main"]