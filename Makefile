# ==============================================================================
# 变量定义
# ==============================================================================

# 自动检测系统架构 (Mac M1/M2 = arm64, Intel = amd64)
# 如果是 Windows Git Bash，可能需要手动指定，例如 make build ARCH=amd64
ARCH ?= $(shell uname -m)
OS ?= $(shell uname -s | tr A-Z a-z)

# 项目基本信息
BINARY = main
OUTPUT = $(BINARY)
# 指向你的 main.go 路径
MAIN_FILE = ./cmd/server/main.go

# Docker 相关
IMAGE_NAME = agentgo-server:v1.0
CONTAINER_NAME = agentgo-app
# 与 deploy/docker-compose.yml 中定义的网络名一致
NETWORK_NAME = agentgo_network

# Go 命令
GO = go
GO_BUILD = $(GO) build
# -w -s 去掉调试信息，减小体积
GO_LDFLAGS = -w -s

.PHONY: all build run clean help

default: run

# ==============================================================================
# 本地开发 (Local Development)
# ==============================================================================

# 1. 直接运行 (最常用)
run:
	@echo ">> Running application locally..."
	$(GO) run $(MAIN_FILE)

# 2. 本地编译
build:
	@echo ">> Building binary for $(OS)/$(ARCH)..."
	$(GO_BUILD) -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT) $(MAIN_FILE)
	@echo ">> Build success: ./$(OUTPUT)"

# 3. 清理编译文件
clean:
	@echo ">> Cleaning..."
	rm -f $(OUTPUT)
	rm -f coverage.out

# ==============================================================================
# 环境依赖 (Infrastructure)
# ==============================================================================

# 启动 MySQL 和 Redis (基于 docker-compose)
env-up:
	@echo ">> Starting dependency containers (MySQL, Redis)..."
	docker compose -f deploy/docker-compose.yml up -d

# 关闭并删除 MySQL 和 Redis 容器
env-down:
	@echo ">> Stopping dependency containers..."
	docker compose -f deploy/docker-compose.yml down

# ==============================================================================
# Docker 部署 (Deployment)
# ==============================================================================

# 1. 构建镜像 (使用你的 Dockerfile)
docker-build:
	@echo ">> Building Docker image: $(IMAGE_NAME)"
	docker build -t $(IMAGE_NAME) .

# 2. 运行容器
# 注意：这里使用了 --network 参数，确保 App 能连上 MySQL
docker-run:
	@echo ">> Running container: $(CONTAINER_NAME)"
	# 如果容器已存在，先删除
	-docker rm -f $(CONTAINER_NAME)
	docker run \
	-d \
	--name $(CONTAINER_NAME) \
	--network $(NETWORK_NAME) \
	-p 8080:8080 \
	$(IMAGE_NAME)
	@echo ">> Container started at http://localhost:8080"

# 3. 停止并删除容器和镜像
docker-clean:
	-docker stop $(CONTAINER_NAME)
	-docker rm $(CONTAINER_NAME)
	-docker rmi $(IMAGE_NAME)

# ==============================================================================
# 帮助信息
# ==============================================================================
help:
	@echo "Make 命令说明:"
	@echo "  make run          - 本地直接运行代码 (go run)"
	@echo "  make build        - 本地编译二进制文件"
	@echo "  make env-up       - 启动 MySQL/Redis 环境"
	@echo "  make env-down     - 关闭 MySQL/Redis 环境"
	@echo "  make docker-build - 构建项目的 Docker 镜像"
	@echo "  make docker-run   - 运行项目容器 (需先运行 make env-up)"