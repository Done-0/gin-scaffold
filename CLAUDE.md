# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

gin-scaffold is a production-ready Go web service scaffold built on Gin framework with enterprise-grade features including database support (PostgreSQL, MySQL, SQLite), Redis caching, Kafka messaging, SSE streaming, i18n, and AI provider integrations (OpenAI, Gemini).

## Development Commands

### Running the Application

```bash
# Local development
go run main.go

# Production mode
ENV=prod go run main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test ./internal/utils/vo/vo_utils_test.go

# Run tests with verbose output
go test -v ./...
```

### Code Quality

```bash
# Format code
gofmt -w .

# Auto-import management
goimports -w .

# Static analysis
go vet ./...

# Generate Wire dependency injection code
cd pkg/wire && wire
```

## Architecture Overview

### Dependency Injection with Wire

The application uses Google Wire for compile-time dependency injection. All dependencies are wired together in `pkg/wire/`:
- `wire.go` defines the dependency graph using `//go:build wireinject`
- `wire_gen.go` is auto-generated (never edit manually)
- `providers.go` contains all provider functions
- The `Container` struct holds all initialized dependencies (managers, controllers, services)

To add new dependencies:
1. Add provider function to `providers.go`
2. Add to `AllProviders` wire set
3. Add field to `Container` struct in `wire.go`
4. Run `wire` to regenerate `wire_gen.go`

### Application Lifecycle

The application follows this initialization sequence (see `cmd/gin_server.go:25`):
1. Load configuration from `configs/config.{local|prod}.yml` based on ENV
2. Wire all dependencies via `wire.NewContainer()`
3. Initialize managers in order: Logger → Database → Redis
4. Setup Gin engine with middleware and routes
5. Start HTTP server with graceful shutdown support

All managers implement `Initialize()` and `Close()` methods and are deferred for proper cleanup.

### Module Organization

The codebase follows a domain-driven structure:

**Infrastructure Layer** (`internal/`):
- `db/` - Database management with GORM (auto-migration, multi-DB support)
- `redis/` - Redis client with connection pooling
- `logger/` - Logrus-based structured logging with rotation
- `queue/` - Kafka producer/consumer
- `sse/` - Server-Sent Events manager
- `i18n/` - Internationalization with go-i18n
- `ai/` - Multi-provider AI service (OpenAI/Gemini) with rate limiting and prompt management

**Domain Layer** (`internal/`):
- `model/` - GORM models with base model pattern
- `types/` - Domain types, constants, error codes
- `utils/` - Shared utilities (errorx, validator, snowflake, etc.)

**Application Layer** (`pkg/`):
- `serve/controller/` - HTTP controllers with DTO validation
- `serve/service/` - Business logic layer
- `serve/mapper/` - Data access layer (optional, use GORM directly if simple)
- `vo/` - View objects for API responses
- `router/` - Route registration organized by modules

### Three-Layer Architecture Pattern

Follow this structure when adding new features:
```
pkg/serve/
├── controller/
│   └── {module}/
│       ├── dto/           # Request DTOs with validation tags
│       └── {module}.go    # HTTP handlers
├── service/
│   └── {module}/
│       ├── {module}.go    # Service interface
│       └── impl/
│           └── {module}.go # Service implementation
└── mapper/
    └── {module}/
        ├── {module}.go    # Mapper interface
        └── impl/
            └── {module}.go # Data access implementation
```

### Configuration Management

Configuration uses Viper with hot-reload support:
- `configs/config.local.yml` - Local development
- `configs/config.prod.yml` - Production
- Config changes are automatically detected and logged
- Access config via `configs.GetConfig()` with RWMutex protection
- Use `UpdateField()` to programmatically modify config

### Error Handling System

The custom error system (`internal/utils/errorx/`) provides:
- Structured errors with status codes and i18n messages
- Stack trace capture for debugging
- Template parameter support for dynamic messages
- Error code registry with range allocation (see `internal/types/errno/README.md`)

Register error codes in `internal/types/errno/`:
```go
errorx.Register(10001, "user not found")
```

Create errors:
```go
return errorx.New(10001)
return errorx.New(10002, errorx.KV("field", "email"))
```

Error code ranges:
- 10000-19999: System errors (next: 10009)

### Middleware Stack

Middleware is registered in `internal/middleware/middleware.go:13`:
1. RequestID - Adds unique request ID
2. Logger - Request/response logging
3. Recovery - Panic recovery
4. CORS - Cross-origin support
5. Custom middleware as needed

### Route Organization

Routes are organized by module in `pkg/router/routes/`:
- Each module has its own routes file (e.g., `test_routes.go`)
- Routes are registered by API version (`/api/v1`, `/api/v2`)
- Registration function pattern: `RegisterXxxRoutes(container, v1, v2)`

## Coding Standards

This project follows strict Go coding conventions documented in `docs/coding-standards.en-US.md`. Key points:

### Naming Conventions
- Packages: lowercase, match directory name
- Files: lowercase with underscores (e.g., `user_service.go`)
- Structs/Interfaces: CamelCase, exported if uppercase first letter
- Variables: camelCase, special acronyms like `API`, `ID` keep original case unless first word
- Booleans: must start with `Has`, `Is`, `Can`, `Allow`
- Constants: CamelCase with category prefix (e.g., `MethodGET`)

### Comments (English only)
- Package comment required with: description, `Author: [GitHub]`, `Created: YYYY-MM-DD`
- Struct/Interface: `// [Name], [description]` format
- Struct fields: inline comments aligned
- Functions: description, parameters, returns (optional but recommended)

### Code Style
- Use `gofmt` for formatting (tabs, not spaces)
- Max 120 characters per line
- Opening brace on same line
- Imports grouped: stdlib → third-party → internal
- Early returns for error handling (no else blocks)
- Test files end with `_test.go`, functions start with `Test`

### Import Organization
```go
import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/Done-0/gin-scaffold/internal/logger"
    "github.com/Done-0/gin-scaffold/pkg/vo"
)
```

## AI Provider Integration

The AI module (`internal/ai/`) supports multiple providers with:
- **Load balancing**: Round-robin across API keys
- **Rate limiting**: Per-instance rate limits (e.g., "60/min")
- **Retry logic**: Configurable max retries
- **Prompt templates**: YAML-based templates in `configs/prompts/`
- **Stream support**: SSE-based streaming responses

Configuration in `configs/config.*.yml` under `AI.PROVIDERS`:
- Each provider (openai, gemini) can have multiple instances
- Each instance has its own keys, models, and parameters
- Enable/disable at provider or instance level

## Common Patterns

### Creating a New API Endpoint

1. Define error codes in `internal/types/errno/system.go`
2. Create DTO in `pkg/serve/controller/dto/` with validation tags
3. Create VO in `pkg/vo/` for response
4. Implement service interface and implementation
5. Create controller with dependency injection
6. Register routes in `pkg/router/routes/`
7. Add provider to `pkg/wire/providers.go`
8. Update `Container` in `pkg/wire/wire.go`

### Database Models

Extend `internal/model/base/base.go` for automatic timestamps:
```go
type User struct {
    base.BaseModel
    Username string `gorm:"uniqueIndex;not null"`
}
```

Models auto-migrate on startup via `internal/db/internal/setup.go:15`.

### Working with Redis

Access via `RedisManager` in container:
```go
rdb := container.RedisManager.GetClient()
rdb.Set(ctx, "key", "value", 0)
```

### SSE Streaming

Use `SSEManager` for server-sent events:
```go
stream := container.SSEManager.CreateStream(clientID)
defer container.SSEManager.RemoveStream(clientID)
// Send events via stream
```

## Environment Variables

- `ENV`: Set to `prod`/`production` for production mode, otherwise local config is used
- Config files are selected automatically based on ENV
