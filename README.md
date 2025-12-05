# Gin HTTP 项目模板

基于 Go 语言 Gin 框架封装的企业级 HTTP 项目模板，开箱即用，包含常用功能模块。

## 特性

- 🚀 基于 Gin 框架，高性能 HTTP 服务
- 🔌 灵活的服务架构（支持 HTTP、TCP、WebSocket 等多种服务）
- 📝 完善的日志记录系统（支持日志分级、文件切割）
- 🛡️ 统一的异常处理机制
- 🔐 用户认证授权（JWT Token）
- ⚙️ 灵活的配置管理（支持多环境配置）
- 🔌 常用中间件集成
- 💾 数据库集成（GORM）
- 📦 优雅的项目结构
- 📄 Swagger API 文档
- ⏱️ 优雅关闭与健康检查

## 技术栈

- Go 1.24+
- Gin Web Framework
- GORM（ORM 框架）
- JWT-Go（认证）
- Viper（配置管理）
- Zap（日志框架）
- Redis（缓存）/ 本地缓存
- MySQL/sqllite（数据库）
- Swagger（API 文档）

## 项目结构

```
.
├── cmd                             # 应用程序入口
│   └── server
│       └── main.go                 # 程序入口
├── config                          # 配置相关
│   ├── config.go                   # 配置结构体定义
│   └── config.yaml                 # 默认配置文件
├── data                            # 数据文件 (例如: 数据库文件, 静态数据)
├── docs                            # Swagger API 文档
│   ├── docs.go                     # Swagger 文档
│   ├── swagger.json                # Swagger JSON
│   └── swagger.yaml                # Swagger YAML
├── internal                        # 内部代码，不向外部暴露
│   ├── api                         # API 接口定义和处理
│   │   ├── v1                      # API v1 版本
│   │   │   ├── auth.go             # 认证相关接口
│   │   │   └── user.go             # 用户相关接口
│   │   └── router.go               # 路由注册
│   ├── controller                  # 控制器层，处理 HTTP 请求
│   │   ├── auth.go                 # 认证控制器
│   │   ├── base.go                 # 基础控制器
│   │   ├── dict.go                 # 字典控制器
│   │   └── user.go                 # 用户控制器
│   ├── dao                         # 数据访问对象层 (Data Access Object)
│   │   ├── base.go                 # 基础 DAO
│   │   ├── dict.go                 # 字典 DAO
│   │   ├── space.go                # 空间 DAO
│   │   └── user.go                 # 用户 DAO
│   ├── database                    # 数据库连接和初始化
│   │   └── database.go             # 数据库连接
│   ├── middleware                  # HTTP 中间件
│   ├── model                       # 业务模型定义
│   │   ├── base.go                 # 基础模型（通用字段）
│   │   ├── dict.go                 # 字典模型
│   │   ├── space.go                # 空间模型
│   │   └── user.go                 # 用户模型
│   ├── server                      # 服务启动器
│   │   └── http.go                 # HTTP 服务启动
│   ├── service                     # 业务逻辑服务层
│   │   ├── auth.go                 # 认证服务
│   │   ├── dict.go                 # 字典服务
│   │   ├── service.go              # 服务接口定义
│   │   └── user.go                 # 用户服务
├── logs                            # 日志文件目录
├── pkg                             # 存放可重用代码包 (公共工具、库)
│   ├── app                         # 应用程序通用逻辑
│   │   └── app.go                  # 应用封装
│   ├── cache                       # 缓存模块
│   │   └── cache.go                # 缓存接口
│   ├── common                      # 通用常量和枚举
│   │   └── const.go                # 常量定义
│   ├── config                      # 配置管理 (例如: Viper)
│   │   └── config.go               # 配置结构体定义
│   ├── enum                        # 枚举定义
│   │   └── user.go                 # 用户枚举
│   ├── errcode                     # 错误码定义
│   │   ├── errcode.go              # 错误码
│   │   └── user.go                 # 用户错误码
│   ├── jwt                         # JWT 相关工具
│   │   └── jwt.go                  # JWT 工具
│   ├── logger                      # 日志模块
│   │   ├── logger.go               # 日志封装
│   │   └── zap.go                  # Zap 日志实现
│   ├── request                     # 请求 DTO (Data Transfer Object)
│   │   ├── request.go              # 请求封装
│   │   └── user.go                 # 用户请求
│   ├── response                    # 响应 DTO
│   │   ├── response.go             # 统一响应格式
│   │   └── user.go                 # 用户响应
│   ├── utils                       # 通用工具函数
│   │   └── utils.go                # 工具函数
│   └── validator                   # 参数验证器
│       └── validator.go            # 参数验证器
├── scripts                         # 脚本文件 (构建、部署、数据库迁移等)
│   ├── build.sh                    # 构建脚本
│   └── migrate.sh                  # 迁移脚本
├── test                            # 测试相关
│   ├── api                         # API 测试
│   │   └── base_test.go            # 基础 API 测试
│   └── unit                        # 单元测试
│       └── service_test.go         # 服务单元测试
├── .env.example                    # 环境变量示例
├── .gitignore                      # Git 忽略文件
├── Dockerfile                      # Docker 构建文件
├── docker-compose.yml              # Docker Compose 配置
├── go.mod                          # Go 模块文件
├── go.sum                          # Go 模块校验文件
├── Makefile                        # Make 命令文件
└── README.md                       # 项目说明文档
```

## 快速开始

### 环境要求

- Go 1.24+
- MySQL 8.0+ 或 sqllite
- Redis 6.0+

### 安装

1. 克隆项目

```bash
git clone <repository-url>
cd gin-template
```

2. 安装依赖

```bash
go mod download
```

3. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库、Redis 等信息
```

4. 初始化数据库

```bash
# 执行数据库初始化脚本
mysql -u root -p < scripts/init.sql
```

5. 运行项目

```bash
# 开发环境（默认启动 HTTP 服务）
go run cmd/server/main.go

# 或使用 Make 命令
make run

# 指定环境
ENV=prod go run cmd/server/main.go

# 启动多个服务（如 HTTP + WebSocket）
go run cmd/server/main.go --services=http,websocket
```

### Docker 部署

```bash
# 构建镜像
docker build -t gin-template:latest .

# 使用 docker-compose 启动
docker-compose up -d
```

## 配置说明

### 配置文件结构

```yaml
server:
  mode: debug # 运行模式：debug/release/test
  http:
    enabled: true # 是否启用 HTTP 服务
    port: 8080 # HTTP 服务端口
    read_timeout: 60 # 读取超时时间（秒）
    write_timeout: 60 # 写入超时时间（秒）
  tcp:
    enabled: false # 是否启用 TCP 服务
    port: 9090 # TCP 服务端口
  grpc:
    enabled: false # 是否启用 gRPC 服务
    port: 9091 # gRPC 服务端口
  websocket:
    enabled: false # 是否启用 WebSocket 服务
    port: 8081 # WebSocket 服务端口

database:
  type: mysql # 数据库类型：mysql/postgres
  host: localhost
  port: 3306
  username: root
  password: password
  dbname: gin_template
  max_idle_conns: 10 # 最大空闲连接数
  max_open_conns: 100 # 最大打开连接数
  conn_max_lifetime: 3600 # 连接最大生命周期（秒）

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  pool_size: 10

jwt:
  secret: your-secret-key
  expire_time: 7200 # Token 过期时间（秒）
  refresh_expire_time: 604800 # 刷新 Token 过期时间（秒）

logger:
  level: info # 日志级别：debug/info/warn/error
  format: json # 日志格式：json/text
  output: logs/app.log # 日志输出路径
  max_size: 100 # 单个日志文件最大大小（MB）
  max_backups: 10 # 保留的旧日志文件数量
  max_age: 30 # 日志文件保留天数
  compress: true # 是否压缩旧日志
```

## API 文档

启动项目后，访问 Swagger 文档：

```
http://localhost:8080/swagger/index.html
```

### 主要接口

#### 认证相关

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出
- `POST /api/v1/auth/refresh` - 刷新 Token
- `POST /api/v1/auth/forgot-password` - 忘记密码
- `POST /api/v1/auth/reset-password` - 重置密码

#### 用户相关

- `GET /api/v1/users/profile` - 获取用户信息
- `PUT /api/v1/users/profile` - 更新用户信息
- `PUT /api/v1/users/password` - 修改密码
- `GET /api/v1/users` - 用户列表（管理员）
- `DELETE /api/v1/users/:id` - 删除用户（管理员）

#### 数据字典

- `GET /api/v1/dict/types` - 获取字典类型列表
- `GET /api/v1/dict/data/:type` - 根据类型获取字典数据
- `POST /api/v1/dict` - 创建字典（管理员）
- `PUT /api/v1/dict/:id` - 更新字典（管理员）
- `DELETE /api/v1/dict/:id` - 删除字典（管理员）

#### 健康检查

- `GET /health` - 健康检查接口

## 中间件说明

### 日志中间件

记录每个请求的详细信息，包括请求方法、路径、耗时、状态码等。

### 异常恢复中间件

捕获 panic 异常，防止服务崩溃，并记录错误日志。

### JWT 认证中间件

验证用户身份，保护需要认证的接口。

### 跨域中间件

处理跨域请求，支持 CORS 配置。

### 限流中间件

基于令牌桶算法的接口限流，防止恶意请求。

### 请求 ID 中间件

为每个请求生成唯一 ID，便于日志追踪。

### 参数验证中间件

自动验证请求参数，支持多种验证规则。

## 错误码定义

```go
const (
    SUCCESS           = 200   // 成功
    ERROR             = 500   // 通用错误
    INVALID_PARAMS    = 400   // 参数错误
    UNAUTHORIZED      = 401   // 未授权
    FORBIDDEN         = 403   // 禁止访问
    NOT_FOUND         = 404   // 资源不存在

    // 用户相关错误码 1xxx
    ERROR_USER_NOT_FOUND     = 1001
    ERROR_USER_EXIST         = 1002
    ERROR_PASSWORD_WRONG     = 1003
    ERROR_TOKEN_INVALID      = 1004
    ERROR_TOKEN_EXPIRED      = 1005

    // 数据字典错误码 2xxx
    ERROR_DICT_NOT_FOUND     = 2001
    ERROR_DICT_TYPE_EXIST    = 2002
)
```

## 统一响应格式

```json
{
  "code": 200,
  "msg": "success",
  "data": {},
  "request_id": "abc123",
  "timestamp": 1699999999
}
```

## Make 命令

```bash
make run          # 运行项目
make build        # 编译项目
make test         # 运行测试
make cover        # 测试覆盖率
make lint         # 代码检查
make fmt          # 格式化代码
make swagger      # 生成 Swagger 文档
make docker       # 构建 Docker 镜像
make clean        # 清理编译文件
```

## 开发指南

### 服务架构说明

项目采用灵活的多服务架构，支持同时运行多种类型的服务：

- **HTTP 服务** (`internal/server/http.go`): 基于 Gin 的 RESTful API 服务
- **TCP 服务** (`internal/server/tcp.go`): 原生 TCP Socket 服务，适用于长连接、自定义协议
- **gRPC 服务** (`internal/server/grpc.go`): 高性能 RPC 服务
- **WebSocket 服务** (`internal/server/websocket.go`): 实时双向通信服务

所有服务共享同一套业务逻辑（Service 层）和数据访问层（Repository 层），只是协议层不同。

### 启动特定服务

```go
// 在 cmd/server/main.go 中
func main() {
    // 初始化配置、数据库等

    // 启动 HTTP 服务
    if config.Server.HTTP.Enabled {
        go server.StartHTTP()
    }

    // 启动 TCP 服务
    if config.Server.TCP.Enabled {
        go server.StartTCP()
    }

    // 启动 gRPC 服务
    if config.Server.GRPC.Enabled {
        go server.StartGRPC()
    }

    // 启动 WebSocket 服务
    if config.Server.WebSocket.Enabled {
        go server.StartWebSocket()
    }

    // 优雅关闭
    gracefulShutdown()
}
```

### 添加新的服务类型

如需添加新的服务类型（如 MQTT、消息队列消费者等）：

1. 在 `internal/server` 目录下创建新的服务文件（如 `mqtt.go`）
2. 实现服务的启动和关闭逻辑
3. 在配置文件中添加相应的配置项
4. 在 `cmd/server/main.go` 中注册并启动服务

### 添加新接口

1. 在 `internal/model` 中定义数据模型
2. 在 `pkg/dto/request` 和 `pkg/dto/response` 中定义 DTO
3. 在 `internal/repository` 中实现数据访问层
4. 在 `internal/service` 中实现业务逻辑
5. 在 `internal/api/v1` 中实现 HTTP 处理器
6. 在 `internal/api/router.go` 中注册路由

### 数据库迁移

使用 GORM 的 AutoMigrate 功能：

```go
db.AutoMigrate(&model.User{}, &model.Dict{})
```

### 添加新的中间件

在 `internal/middleware` 目录下创建中间件文件，并在路由中注册。

## 测试

```bash
# 运行所有测试
go test ./...

# 运行单个包的测试
go test ./internal/service

# 查看测试覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 性能优化建议

- 使用连接池管理数据库连接
- 合理使用 Redis 缓存热点数据
- 对频繁查询的接口添加缓存
- 使用分页查询避免大量数据加载
- 对慢查询进行优化
- 使用索引提升查询性能
- 启用 GZIP 压缩

## 安全建议

- 所有密码使用 bcrypt 加密存储
- 敏感信息不要记录到日志
- 使用 HTTPS 传输数据
- 定期更新依赖包
- 对输入参数进行严格验证
- 实施接口限流防止暴力破解
- 定期备份数据库

## 生产部署清单

- [ ] 修改配置文件中的敏感信息
- [ ] 设置 Gin 为 release 模式
- [ ] 配置日志级别为 info 或 warn
- [ ] 启用 HTTPS
- [ ] 配置防火墙规则
- [ ] 设置数据库连接池参数
- [ ] 配置 Redis 持久化
- [ ] 设置合理的超时时间
- [ ] 配置反向代理（Nginx）
- [ ] 配置自动重启（Systemd/Supervisor）
- [ ] 配置监控告警

## 常见问题

### 1. 数据库连接失败

检查配置文件中的数据库连接信息是否正确，确保数据库服务已启动。

### 2. JWT Token 无效

检查 Token 是否已过期，或者 Secret Key 是否正确配置。

### 3. 跨域问题

确保跨域中间件已正确配置，检查允许的来源域名。

### 4. 日志文件过大

配置日志切割参数，设置合理的 max_size 和 max_backups。

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## License

[MIT License](LICENSE)

## 联系方式

- 项目主页: <repository-url>
- 问题反馈: <repository-url>/issues
- 邮箱: your-email@example.com

---

**注意**: 在生产环境使用前，请务必修改所有默认配置和密钥！
