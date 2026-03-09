# AgentGo

轻量用户模块示例，包含 DAO / Service / Controller 分层与简单的注册、登录、登出流程。

## 本地运行

- 依赖：Go 1.25（推荐），可用 MySQL/Redis（非必须用于测试）。
- 配置：生产运行需在 `config/config.yaml` 中配置数据库、Redis 与邮件服务；本次测试不依赖该配置。

## 运行测试

本仓库的测试覆盖：
- DAO：`internal/dao/user_test.go`
- Service：`internal/service/user_test.go`

安装依赖并运行测试：

```bash
go test ./...
```

## 测试设计说明

- DAO 测试：使用 `gorm.io/driver/sqlite` 的内存库，`AutoMigrate` 出 `User` 表，验证 `CreateUser`、`GetUserByEmail`、`CheckUserExist` 行为。
- Service 测试：通过内存 mock（实现 `UserDao` 与 `UserCacheDao` 接口）验证注册/登录的主要分支：
  - 注册成功（验证码正确，用户不存在）
  - 注册失败（用户已存在 / 验证码不匹配）
  - 登录成功、密码错误、用户不存在
- Controller 测试：
  - 使用 Gin 的 `httptest` 驱动 `UserRegister` / `UserLogin` / `UserLogout` 路由
  - 将全局 `db.DB` 指向 sqlite 内存库；将全局 `redis.RDB` 指向 miniredis 模拟服务
  - 预先写入验证码 Key，完整校验注册 → 登录 → 登出流程

## 常见问题

- 如遇到 Go 版本不兼容，请升级到较新的 Go 版本（建议 1.25）。

## 目录结构

详见仓库根目录树，主要目录：
- `internal/api/v1`：Gin 控制器（HTTP 层）
- `internal/service`：业务逻辑层
- `internal/dao`：数据访问层（GORM）
- `internal/common`：MySQL、Redis 初始化与工具
- `pkg`：通用工具（配置、JWT、响应包装器、错误码等）