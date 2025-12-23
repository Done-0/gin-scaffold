# gin-scaffold 脚手架

## 简介

gin-scaffold 是一个现代化的 Go Web 服务端脚手架，基于 Gin 框架，集成主流企业级能力，并在 docs 目录下附带了完善的 Go 语言开发规范，助力团队高效、规范地推进中大型项目开发。

## 技术栈

- **Go 1.25.1**
- **Gin**：Web 框架，路由与中间件支持
- **GORM**：ORM 框架，支持 PostgreSQL、MySQL、SQLite，自动迁移
- **go-redis/v9**：Redis 客户端，连接池与健康检查
- **Logrus** + **file-rotatelogs/lfshook**：结构化日志，分级与轮转
- **Viper** + **fsnotify**：配置管理，支持多环境与热加载
- **Wire**：依赖注入
- **Kafka (sarama)**：消息队列，生产者/消费者
- **SSE (gin-contrib/sse)**：服务端推送
- **go-playground/validator/v10**：参数校验
- **go-i18n/v2**：国际化
- **bwmarrin/snowflake**：雪花 ID
- **中间件**：RequestID、CORS、Gzip、Secure、Recovery、事务、验证码等
- **工具库**：文件处理、模板渲染、速率限制、分布式队列、错误处理、VO、AI（OpenAI/Gemini）、API 示例等

## 快速开始

1. 克隆本仓库并安装依赖

```bash
git clone https://github.com/Done-0/gin-scaffold.git
```

2. 根据 example 文件配置 `configs/config.local.yaml` 或 `configs/config.prod.yaml`

3. 启动服务

```bash
go run main.go
```

## 适用场景

- 企业级后端系统开发
- 微服务架构与分布式服务
- RESTful API 服务
- 高并发与高可用业务场景
- 快速原型及功能验证
- 追求高可维护性、易扩展的 Go 项目

## 主要模块说明

- **数据库模块**：支持 PostgreSQL、MySQL、SQLite，自动建库与迁移，灵活参数配置，统一管理入口，基于 GORM 实现
- **缓存模块**：全局 Redis 客户端，连接池、超时、健康检查，统一管理，适用于缓存和分布式场景
- **日志模块**：统一接口，结构化日志，分级、自动轮转，便于生产环境分析与追踪
- **中间件与工具**：集成请求 ID、CORS、安全、恢复、Gzip、事务、验证码、参数校验、雪花 ID等常用功能，提升安全性与开发效率
- **优雅启动与关闭**：支持信号优雅关闭，自动释放数据库和缓存资源，保障服务稳定性
- **API与测试示例**：内置多组接口和测试用例，便于快速验证各模块功能和业务逻辑

## 架构推荐

### 经典三层架构

```bash
./pkg
├── ./pkg/router
│   ├── ./pkg/router/routes # 路由组
│   │   ├── ./pkg/router/routes/test_routes.go
│   │   └── ./pkg/router/routes/user_routes.go
│   └── ./pkg/router/router.go
├── ./pkg/serve
│   ├── ./pkg/serve/controller # controller 控制层
│   │   └── ./pkg/serve/controller
│   │       ├── ./pkg/serve/controller/dto # dto
│   │       │   ├── ./pkg/serve/controller/dto/user_dto.go
│   │       │   └── ./pkg/serve/controller/dto/test_dto.go
│   │       ├── ./pkg/serve/controller/user_controller.go
│   │       └── ./pkg/serve/controller/test_controller.go
│   ├── ./pkg/serve/mapper # mapper 层
│   │   └── ./pkg/serve/mapper
│   │       ├── ./pkg/serve/mapper/impl
│   │       │   ├── ./pkg/serve/mapper/impl/user_mapper_impl.go
│   │       │   └── ./pkg/serve/mapper/impl/test_mapper_impl.go
│   │       ├── ./pkg/serve/mapper/test_mapper.go
│   │       └── ./pkg/serve/mapper/user_mapper.go
│   └── ./pkg/serve/service # service 服务层
│       └── ./pkg/serve/service
│           ├── ./pkg/serve/service/impl
│           │   ├── ./pkg/serve/service/impl/user_service_impl.go
│           │   └── ./pkg/serve/service/impl/test_service_impl.go
│           ├── ./pkg/serve/service/user_service.go
│           └── ./pkg/serve/service/test_service.go
└── ./pkg/vo
    └── ./pkg/vo # vo
        ├── ./pkg/vo/user_vo.go
        └── ./pkg/vo/test_vo.go
```

## 贡献

欢迎 issue 和 PR！

## 许可证

本项目采用 Apache License 2.0 协议开源
