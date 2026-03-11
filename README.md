<div align="center">

# AgentGo

面向“企业级多智能体（Multi-Agents）研发自动化”的 Go 项目脚手架，内置分层架构、用户模块与可扩展的 LLM 适配层，支持流式对话上下文与持久化，为后续 RAG 能力与多智能体协同（后端/测试/环境搭建/自动化部署/前端）铺路。

</div>


**目录**
- 项目介绍
- 项目优势
- 愿景与路线图
- 架构总览
- 核心模块
- 快速开始（安装/配置/测试）
- 使用示例（同步/流式）
- 配置说明
- 扩展指南（LLM 适配/RAG/多智能体）
- 开发规范（目录结构/依赖/测试）
- 常见问题与注意事项
 - 错误码
- 贡献、许可与安全

---

## 项目介绍

AgentGo 旨在构建一个可落地的“多智能体研发协作平台”后端基座：
- 统一的 LLM 适配与流式输出能力
- 会话/消息模型与持久化，支撑多轮对话与上下文管理
- 渐进式引入 RAG，增强知识检索与事实一致性
- 逐步引入 Multi-Agents，覆盖研发全流程（后端 → 测试 → 环境搭建 → 持续交付 → 前端）

当前仓库已包含用户模块、LLM 接入与基础会话能力，后续将持续迭代更完整的自动化链路。

## 项目优势

- 无状态架构：会话上下文不常驻内存，不依赖进程内全局管理器，天然支持多副本、负载均衡与弹性扩缩容；请求到达任意实例都可通过数据库快速恢复上下文，避免状态漂移与 OOM 风险。
- `Helper` 请求级生命周期：保留并发安全（`RWMutex`）、上下文组装与 `saveFunc` 持久化回调的优点，将其作用域限定在单次请求中，完成流式/同步生成后自然释放内存，降低资源占用与内存泄漏风险。
- 工厂与自动注册：适配器在 `init()` 中自注册，`factory.CreateModel()` 统一构造模型实例，扩展新模型无需侵入式改动，解耦良好。
- 持久化策略可插拔：默认直写 MySQL，链路简单、延迟低；若未来需要上 MQ（如 RabbitMQ/Kafka）实现异步存储，可仅替换 `saveFunc`，核心业务与适配层零改动。
- 流式实践完善：SSE 回包遵循规范，首包可返回 `session_id` 便于前端落盘，增量 `data:` 连续下发，结尾 `event: end`/失败 `event: error`，并补充了跨域与代理缓冲相关响应头。
- 配置与部署友好：`LLMConfig` 配置化选择模型与网关，支持 `APP_ENV` 区分环境；提供 Dockerfile/Compose/Makefile，一键构建与启动。
- 工程化测试：DAO/Service 覆盖核心路径，易于持续演进；接口错误码清晰分层，便于联调与观测。

前端集成建议（类型与参数）
- `session_id`：当前为 `uint` 类型，请前端以数值提交；如因历史原因使用字符串，可在接入层转换或在服务端增加兼容解码（将字符串数值安全转为 `uint`）。
- `model_type`：需与工厂注册名保持一致（如 `openai`、`ollama`）。
- SSE：浏览器端建议使用 `EventSource`，并处理 `event: error` 与 `event: end` 两类事件；跨域/代理需与网关策略匹配放行。

## 愿景与路线图

阶段 1：对话与上下文
- LLM 统一适配（OpenAI/Ollama），支持同步与流式
- 会话与消息持久化，便捷拉取历史上下文

阶段 2：RAG（Retrieval Augmented Generation）
- 文档摄取与切分、Embedding 生成
- 向量检索召回、重排与答案生成
- 可插拔的向量库接口（本地/云端）

阶段 3：Multi-Agents 协作
- 任务编排与 Agent 能力角色化（后端/测试/环境/DevOps/前端）
- 工程上下文缓存与工具调用（代码分析、运行、部署）
- 端到端流水线（代码 → 测试 → 构建 → 部署 → 验收）

## 架构总览

- Web/API：Gin（位于 internal/api/v1）
- Service：业务编排（internal/service）
- DAO：数据存取（internal/dao，GORM）
- Common：基础设施（internal/common/mysql、redis）
- LLM：统一适配与会话辅助（internal/llm）
- 包工具：配置、错误码、JWT 等（pkg/*）

## 核心模块

### LLM（internal/llm）

提供统一大模型接口、工厂注册与适配器：
- 核心接口与转换：[internal/llm/model.go](internal/llm/model.go)
  - `Model`：`GenerateResponse`（同步）、`StreamResponse`（流式）、`GetModelType`
  - `toEinoMessages`：适配 CloudWeGo EINO 的消息结构
- 适配器：OpenAI、Ollama
  - [internal/llm/adapter_openai.go](internal/llm/adapter_openai.go)
  - [internal/llm/adapter_ollama.go](internal/llm/adapter_ollama.go)
- 工厂与注册：[internal/llm/factory.go](internal/llm/factory.go)
- 会话辅助器：[internal/llm/helper.go](internal/llm/helper.go)

### 会话与消息（Session/Message）

- 模型：
  - 会话：[internal/model/session.go](internal/model/session.go)
  - 消息：[internal/model/message.go](internal/model/message.go)
- DAO：
  - SessionDao（创建/查询/删除）：[internal/dao/session.go](internal/dao/session.go)
  - MessageDao（写入/按会话拉取历史）：[internal/dao/message.go](internal/dao/message.go)
- Service：
  - 会话与对话流式封装：[internal/service/session.go](internal/service/session.go)

## 快速开始

### 先决条件
- Go 1.25+
- 可选：MySQL、Redis（测试可不依赖）

### 安装与配置
1. 拉取依赖
```bash
go mod tidy
```
2. 初始化配置（必要时修改）
- 基础配置文件：[config/config.yaml](config/config.yaml)
- 配置结构体：[pkg/conf/config.go](pkg/conf/config.go)

3. 加载配置（你的程序入口需调用一次）
```go
conf.Init()
```

### 构建与测试
```bash
go build ./...
go test ./...
```

涵盖测试（示例）：
- DAO：
  - 用户：[internal/dao/user_test.go](internal/dao/user_test.go)
  - 会话：[internal/dao/session_test.go](internal/dao/session_test.go)
  - 消息：[internal/dao/message_test.go](internal/dao/message_test.go)
- Service：
  - 用户：[internal/service/user_test.go](internal/service/user_test.go)
  - 会话：[internal/service/session_test.go](internal/service/session_test.go)

## 使用示例

### 使用 Helper 同步生成
```go
ctx := context.Background()
m, _ := llm.CreateModel(ctx, conf.Config.LLM.Type, &conf.Config.LLM)
h := llm.NewHelper(m, 123) // 传入 sessionID
h.SetSaveFunc(func(msg *model.Message) error { /* 持久化 */ return nil })
reply, _ := h.GenerateResponse(ctx, "alice", "你好，给我一个项目结构建议？")
fmt.Println(reply.Content)
```

### 流式输出（SSE/WebSocket 可复用）
```go
_, _ = h.StreamResponse(ctx, "alice", "逐步解释你的设计思路", func(delta string){
    // 将 delta 推送给前端，如 SSE/WS
})
```

### 会话 Service 典型调用
```go
svc := service.NewSessionService(dao.NewSessionDao(db), dao.NewMessageDao(db))
resp, code := svc.CreateSession(ctx, &types.CreateSessionRequest{Username: "alice", Title: "后端设计讨论"})
_ = code
sid := resp.(*types.CreateSessionResponse).SessionID
_, _ = svc.StreamChat(ctx, &types.StreamChatRequest{Username: "alice", SessionID: sid, Question: "给出一个 Gin SSE 样例", ModelType: conf.Config.LLM.Type}, func(delta string){})
```

## API 接口

- 路由初始化：[internal/api/router.go](internal/api/router.go)
  - 前缀：/api/v1
  - 用户：/user/register、/user/login、/user/logout（见 [internal/api/v1/user.go](internal/api/v1/user.go)）
  - 会话：/session/create、/session/history、/session/stream（见 [internal/api/v1/session.go](internal/api/v1/session.go)）

### 统一响应封装

- 定义：[pkg/ctl/ctl.go](pkg/ctl/ctl.go)
- 字段：
  - code：业务错误码（见“错误码”）
  - msg：人类可读信息
  - data：具体数据

示例（成功）：
```json
{
  "code": 200,
  "msg": "ok",
  "data": {"session_id": 123}
}
```

### 用户接口

- 注册：POST /api/v1/user/register
  - 请求体：`UserRegisterRequest`（见 [internal/types/user.go](internal/types/user.go)）
  - 示例：
    ```bash
    curl -sS -X POST http://localhost:8080/api/v1/user/register \
      -H 'Content-Type: application/json' \
      -d '{"email":"u@example.com","captcha":"123456","password":"abc123"}'
    ```

- 登录：POST /api/v1/user/login
  - 请求体：`UserLoginRequest`
  - 示例：
    ```bash
    curl -sS -X POST http://localhost:8080/api/v1/user/login \
      -H 'Content-Type: application/json' \
      -d '{"email":"u@example.com","password":"abc123"}'
    ```

- 登出：POST /api/v1/user/logout
  - Header：Authorization: Bearer <token>
  - 示例：
    ```bash
    curl -sS -X POST http://localhost:8080/api/v1/user/logout \
      -H 'Authorization: Bearer <token>'
    ```

### 会话接口

- 创建：POST /api/v1/session/create
  - 请求体：`CreateSessionRequest{ username, title }`
  - 响应体：`CreateSessionResponse{ session_id(uint) }`
  - 示例：
    ```bash
    curl -sS -X POST http://localhost:8080/api/v1/session/create \
      -H 'Content-Type: application/json' \
      -d '{"username":"alice","title":"后端设计讨论"}'
    ```

- 历史：POST /api/v1/session/history
  - 请求体：`GetHistoryRequest{ session_id(uint) }`
  - 响应体：`GetHistoryResponse{ history: [{is_user, content}, ...] }`
  - 示例：
    ```bash
    curl -sS -X POST http://localhost:8080/api/v1/session/history \
      -H 'Content-Type: application/json' \
      -d '{"session_id":123}'
    ```

- 流式：POST /api/v1/session/stream（SSE）
  - 请求体：`StreamChatRequest{ username, session_id, question, model_type }`
  - 响应头：`Content-Type: text/event-stream`；多条 `data: <chunk>`；结束事件 `event: end`，失败 `event: error`
  - 示例：
    ```bash
    curl -N -X POST http://localhost:8080/api/v1/session/stream \
      -H 'Content-Type: application/json' \
      -d '{"username":"alice","session_id":123,"question":"你好","model_type":"openai"}'
    ```

## 配置说明

OpenAI：
```yaml
llm:
  type: "openai"
  api_key: "sk-xxxxxxxxxxxxxxxx"
  base_url: "https://api.openai.com/v1"
  model_name: "gpt-4o"
```

Ollama：
```yaml
llm:
  type: "ollama"
  base_url: "http://localhost:11434/v1"
  model_name: "qwen2.5:7b"
```

`APP_ENV=prod` 时，将加载 `config/config.prod.yaml`（见 [pkg/conf/config.go](pkg/conf/config.go)）。

## 扩展指南

### 扩展 LLM 适配器
实现 `Model` 接口，并在 `init()` 中通过 `Register()` 注册，按需对接 `eino-ext` 生态或自研 HTTP 客户端。

### 引入 RAG（规划）
- 文档管道：采集 → 切分 → 向量化（Embedding） → 入库
- 检索阶段：向量召回 → 候选重排 → 上下文拼接
- 生成阶段：结合检索上下文生成答案，并可选引用出处
- 抽象接口：`VectorStore`、`Embedder`、`DocStore`，支持替换实现（本地/云端）

### Multi-Agents（规划）
- 角色与工具：后端（代码生成/修改）、测试（用例/单测/覆盖率）、环境（Docker/Compose/K8s 清单）、DevOps（CI/CD）、前端（页面/组件）
- 能力路由：根据任务类型选择 Agent Pipeline（链/图/小型 Orchestrator）
- 运行模式：同步（短链路）/ 异步（持久化状态，任务队列）

## 开发规范

### 目录结构（关键路径）
- internal/api/v1：控制器（HTTP）
- internal/service：业务逻辑
- internal/dao：数据访问（GORM）
- internal/common：MySQL/Redis 初始化
- internal/llm：大模型接口与会话辅助
- pkg：通用能力（配置、错误码、JWT、工具）

### 依赖管理
- 依赖声明：go.mod（Go 1.25）
- 推荐命令：
```bash
go mod tidy
go build ./...
go test ./...
```

### 测试
- 单元测试样例：
  - DAO 测试：[internal/dao/user_test.go](internal/dao/user_test.go)
  - Service 测试：[internal/service/user_test.go](internal/service/user_test.go)
- 建议：优先对修改点与新模块补齐单测；对外部依赖（DB/Redis/LLM）使用替身或内存实现

## 常见问题与注意事项
- `SessionID` 在 API 与模型层统一为 `uint`
- `context.Context` 贯穿 DAO（含 `CreateMessage`），便于设置超时/取消
- OpenAI 需有效 `api_key` 与可达 `base_url`；Ollama 需本地服务与正确的 `model_name`

## 错误码

定义位置：[pkg/e/code.go](pkg/e/code.go)、[pkg/e/msg.go](pkg/e/msg.go)

通用/系统级（100xx）

| Code  | Key                          | Message     |
|------:|------------------------------|-------------|
| 10001 | ERROR_AUTH_CHECK_TOKEN_FAIL  | Token鉴权失败 |
| 10002 | ERROR_AUTH_CHECK_TOKEN_TIMEOUT | Token已超时   |

用户模块（200xx）

| Code  | Key                   | Message     |
|------:|-----------------------|-------------|
| 20001 | ERROR_USER_NOT_EXIST  | 用户不存在     |
| 20002 | ERROR_USER_EXIST      | 用户名已存在   |
| 20003 | ERROR_USER_WRONG_PWD  | 密码错误       |
| 20004 | ERROR_INVALID_CAPTCHA | 验证码错误     |
| 20005 | ERROR_SEND_EMAIL      | 发送邮件失败   |

会话/对话/LLM（300xx）

| Code  | Key                         | Message     |
|------:|-----------------------------|-------------|
| 30001 | ERROR_SESSION_CREATE_FAIL   | 会话创建失败   |
| 30002 | ERROR_HISTORY_LOAD_FAIL     | 历史记录加载失败 |
| 30003 | ERROR_LLM_CREATE_FAIL       | 模型初始化失败 |
| 30004 | ERROR_STREAM_RESPONSE_FAIL  | 流式响应失败   |
| 30005 | ERROR_INVALID_MODEL_TYPE    | 无效的模型类型 |

## 贡献、许可与安全
- 贡献：欢迎提交 Issue/PR，建议遵循规范化的分支与提交信息
- 许可：待补充（建议 MIT/Apache-2.0）
- 安全：如涉及敏感问题，请通过私信渠道联系维护者
