# gin-scaffold 脚手架

## 简介

gin-scaffold 是一个现代化的 Go Web 服务端脚手架，基于 Gin 框架，集成主流企业级能力，并在 docs 目录下附带了完善的 Go 语言开发规范，助力团队高效、规范地推进中大型项目开发。

## 技术栈

- **Go 1.23+**
- **Gin**：高性能 Web 框架
- **GORM**：ORM 框架，支持 PostgreSQL、MySQL、SQLite
- **Redis**：go-redis/v9 客户端
- **Logrus**：结构化日志与日志轮转
- **Viper**：配置管理与热加载
- **air**：热重载工具，使用 `air -c ./configs/.air.toml`
- **中间件**：requestid、CORS、Gzip、Secure、Recovery 等
- **工具库**：雪花 ID、参数校验、数据库事务等

## 快速开始

1. 克隆本仓库并安装依赖 `git clone https://github.com/Done-0/gin-scaffold.git`
2. 配置 `configs/configs.yaml`
3. 启动服务 `go run main.go`

## 适用场景

- 企业级后端服务
- 微服务/RESTful API
- 快速原型开发
- 需要高可维护性和扩展性的 Go 项目

## 主要模块说明

- **数据库组件**：支持 PostgreSQL、MySQL、SQLite，自动建库、自动迁移，连接参数灵活配置，基于 GORM 实现 ORM。
- **Redis 组件**：提供全局 Redis 客户端，支持连接池、超时、健康检查，便于缓存和分布式场景。
- **日志组件**：统一日志接口，支持 JSON 格式、日志分级、自动轮转，便于生产环境日志分析。
- **中间件与工具**：包含请求 ID、CORS、安全、恢复、Gzip、数据库事务、验证码、数据校验、雪花 ID 等常用中间件和工具函数。
- **优雅启动/关闭**：支持 Ctrl+C 或 kill 信号优雅关闭服务，自动释放数据库和缓存资源。
- **API 示例**：内置多组测试接口，便于快速验证各模块功能。

## 架构推荐

### 经典三层架构

```bash
./pkg
│   ├── ./pkg/router
│   │   ├── ./pkg/router/routes # 路由组
│   │   │   ├── ./pkg/router/routes/test.go
│   │   │   └── ./pkg/router/routes/user.go
│   │   └── ./pkg/router/router.go
│   ├── ./pkg/serve
│   │   ├── ./pkg/serve/controller # controller 控制层
│   │   │   ├── ./pkg/serve/controller/test
│   │   │   │   └── ./pkg/serve/controller/test/test.go
│   │   │   └── ./pkg/serve/controller/user
│   │   │       ├── ./pkg/serve/controller/user/dto # dto
│   │   │       │   └── ./pkg/serve/controller/user/dto/user.go
│   │   │       └── ./pkg/serve/controller/user/user.go
│   │   ├── ./pkg/serve/mapper # mapper 层
│   │   │   └── ./pkg/serve/mapper/user
│   │   │       ├── ./pkg/serve/mapper/user/impl
│   │   │       │   └── ./pkg/serve/mapper/user/impl/user.go
│   │   │       └── ./pkg/serve/mapper/user/user.go
│   │   └── ./pkg/serve/service
│   │       └── ./pkg/serve/service/user # service 服务层
│   │           ├── ./pkg/serve/service/user/impl
│   │           │   └── ./pkg/serve/service/user/impl/user.go
│   │           └── ./pkg/serve/service/user/user.go
│   └── ./pkg/vo
│       ├── ./pkg/vo/user # vo
│       │   └── ./pkg/vo/user/user.go
│       └── ./pkg/vo/result.go
```

## 贡献

欢迎 issue 和 PR！
