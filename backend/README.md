# Artisan Coder - Backend API

Phase 2: 用户认证系统后端 API

## 技术栈

- **Go 1.21+** - 后端语言
- **Gin** - HTTP Web 框架
- **GORM** - ORM 库
- **PostgreSQL** - 关系型数据库
- **JWT-Go** - JWT Token 认证
- **FX** - 依赖注入框架
- **Viper** - 配置管理

## 项目结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # 应用入口
├── internal/
│   ├── app/
│   │   └── app.go                  # FX 模块组装
│   ├── config/
│   │   └── config.go              # 配置管理
│   ├── database/
│   │   └── database.go            # 数据库模块
│   ├── models/
│   │   └── user.go                # User 模型
│   ├── repository/
│   │   └── user_repository.go     # 数据访问层
│   ├── handler/
│   │   └── auth_handler.go        # HTTP 处理器
│   ├── middleware/
│   │   ├── cors.go                # CORS 中间件
│   │   ├── auth.go                # JWT 认证中间件
│   │   └── logger.go              # 日志中间件
│   ├── service/
│   │   └── auth_service.go        # 业务逻辑层
│   ├── router/
│   │   └── router.go              # 路由模块
│   └── server/
│       └── server.go              # HTTP 服务器模块
├── pkg/
│   ├── jwt/
│   │   └── jwt.go                 # JWT 工具
│   ├── password/
│   │   └── password.go            # 密码加密工具
│   └── response/
│       └── response.go            # 统一响应格式
├── configs/
│   ├── config.development.yaml    # 开发环境配置
│   └── config.production.yaml     # 生产环境配置
├── migrations/                    # 数据库迁移文件
└── go.mod
```

## 快速开始

### 前置要求

- Go 1.21+
- PostgreSQL 15+

### 安装依赖

```bash
cd backend
go mod download
```

### 配置

编辑 `configs/config.development.yaml` 配置数据库连接：

```yaml
database:
  host: "localhost"
  port: "5433"
  user: "artisan"
  password: "artisan_password_change_me"
  dbName: "artisan_coder"
  sslMode: "disable"
```

### 运行

```bash
# 构建
go build -o bin/server ./cmd/server

# 运行
./bin/server
```

服务器将在 `http://localhost:8080` 启动。

## API 接口

### 基础路径

所有 API 路径以 `/api` 为前缀。

### 接口列表

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/auth/register | 用户注册 | 否 |
| POST | /api/auth/login | 用户登录 | 否 |
| POST | /api/auth/logout | 用户登出 | 是 |
| POST | /api/auth/refresh | 刷新 Token | 否 (使用 Refresh Token) |
| GET | /api/auth/me | 获取当前用户 | 是 |

### 请求/响应格式

#### 用户注册

**请求**: `POST /api/auth/register`

```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "confirmPassword": "password123"
}
```

**响应**: `201 Created`

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "johndoe",
      "email": "john@example.com",
      "createdAt": "2025-12-26T10:30:00Z",
      "updatedAt": "2025-12-26T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

#### 用户登录

**请求**: `POST /api/auth/login`

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**响应**: `200 OK`

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user": { ... },
    "token": "...",
    "refreshToken": "..."
  }
}
```

#### 获取当前用户

**请求**: `GET /api/auth/me`

Headers:
```
Authorization: Bearer {token}
```

**响应**: `200 OK`

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "...",
    "username": "johndoe",
    "email": "john@example.com",
    "createdAt": "...",
    "updatedAt": "..."
  }
}
```

### 错误响应

所有错误响应格式统一：

```json
{
  "code": 400,
  "message": "错误消息",
  "data": null
}
```

| 响应码 | HTTP 状态 | 说明 |
|--------|----------|------|
| 0 | 200/201 | 成功 |
| 400 | 400 | 请求参数错误 |
| 401 | 401 | 未授权 |
| 404 | 404 | 资源不存在 |
| 409 | 409 | 资源冲突 |
| 500 | 500 | 服务器内部错误 |

## 环境变量

可以通过环境变量覆盖配置文件：

```bash
# 应用环境
export APP_ENV=development

# 服务器配置
export SERVER_PORT=8080

# 数据库配置
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=artisan
export DB_PASSWORD=artisan123
export DB_NAME=artisan_coder

# JWT 配置
export JWT_SECRET=your-secret-key

# CORS 配置
export FRONTEND_URL=http://localhost:5173
```

## Docker 部署

### 构建镜像

```bash
docker build -t artisan-backend ./backend
```

### 使用 Docker Compose

在项目根目录运行：

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f backend

# 停止服务
docker-compose down
```

## 测试

### 使用 curl 测试

```bash
# 注册
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123","confirmPassword":"password123"}'

# 登录
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# 获取当前用户
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 开发说明

### 添加新的 API 端点

1. 在 `internal/handler/` 中添加处理函数
2. 在 `internal/router/router.go` 中注册路由
3. 在 `internal/service/` 中添加业务逻辑（如需要）
4. 在 `internal/repository/` 中添加数据访问函数（如需要）

### 数据库迁移

在生产环境，使用 `golang-migrate` 运行迁移：

```bash
# 创建迁移
migrate create -ext sql -dir migrations -seq create_new_table

# 执行迁移
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# 回滚迁移
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down 1
```

## 许可证

MIT
