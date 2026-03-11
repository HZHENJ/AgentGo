# 第一阶段：构建 (Builder)
FROM golang:1.25 AS builder

# 设置环境变量，开启 Go Modules，配置国内代理
#ENV GO111MODULE=on \
#    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

# 1. 先拷贝 go.mod 和 go.sum 下载依赖 (利用 Docker 缓存机制)
COPY go.mod go.sum ./
RUN go mod download

# 2. 拷贝源代码
COPY . .

# 3. 编译
# 注意：CGO_ENABLED=0 是为了生成静态链接的可执行文件
# 注意：这里去掉了 GOARCH=arm64，默认使用宿主机架构，或者在 makefile 里控制
# 构建路径指向 cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/server/main.go

# 4. 准备发布目录
WORKDIR /app/publish
RUN mkdir config && \
    cp /app/main . && \
    cp /app/config/config.prod.yaml ./config/config.yaml

# 第二阶段：运行 (Runner)
# 使用 alpine 比 busybox 更常用，且支持包管理 (apk) 方便排查问题
FROM alpine:latest

WORKDIR /app

# 设置时区 (很多 bug 都是因为时区不对造成的)
RUN apk --no-cache add ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 从构建阶段拷贝编译好的文件和配置
COPY --from=builder /app/publish .

# 环境变量
ENV GIN_MODE=release
ENV SERVICE_PORT=:8080

EXPOSE 8080

ENTRYPOINT ["./main"]