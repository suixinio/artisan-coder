# Phase 2: API 后端 - 用户认证系统

> **目标**: 实现用户认证相关的后端 API，包括注册、登录、登出、Token 刷新等功能，与 Phase 1 前端完全对接
>
> **技术栈**: Go + Gin + GORM + PostgreSQL + JWT + Docker
>
> **状态**: 🔄 待开始
> **预计工时**: 5-7 天

---

## 目录

- [1. 项目概述](#1-项目概述)
- [2. 技术选型](#2-技术选型)
- [3. 项目结构](#3-项目结构)
- [4. 数据库设计](#4-数据库设计)
- [5. API 接口设计](#5-api-接口设计)
- [6. 核心模块实现](#6-核心模块实现)
- [7. 开发任务分解](#7-开发任务分解)
- [8. 测试计划](#8-测试计划)
- [9. 部署配置](#9-部署配置)
- [10. 前后端联调](#10-前后端联调)
- [11. 附录](#11-附录)

---

## 1. 项目概述

### 1.1 背景

Phase 1 已完成前端登录注册界面和 Dashboard，使用 Mock 数据进行开发。Phase 2 将实现真实的后端 API，提供用户认证功能。

### 1.2 目标

- [x] 实现用户注册 API
- [x] 实现用户登录 API（JWT Token 认证）
- [x] 实现 Token 刷新机制
- [x] 实现用户登出 API
- [x] 实现获取当前用户信息 API
- [x] 集成 PostgreSQL 数据持久化
- [x] 实现 CORS 和认证中间件
- [x] 提供 Docker 部署方案

### 1.3 范围

**本阶段包含**:
- 用户认证系统（注册、登录、登出、Token 刷新）
- 用户基本信息管理
- 数据库设计与迁移
- JWT Token 生成与验证
- 密码加密存储（bcrypt）

**不包含**:
- 项目管理功能（Phase 3）
- 会话管理（Phase 3）
- WebSocket 通信（Phase 4）
- 文件操作（Phase 4）

---

## 2. 技术选型

### 2.1 后端框架

| 技术 | 版本 | 用途 |
|------|------|------|
| **Go** | 1.21+ | 后端开发语言 |
| **Gin** | 1.9+ | HTTP Web 框架 |
| **GORM** | 1.25+ | ORM 库 |
| **PostgreSQL** | 15+ | 关系型数据库 |
| **JWT-Go** | 5.0+ | JWT Token 生成与验证 |

### 2.2 开发工具

| 工具 | 用途 |
|------|------|
| **golang-migrate** | 数据库迁移工具 |
| **air** | 热重载开发工具 |
| **swag** | API 文档生成（Swagger） |
| **testify** | 单元测试框架 |

### 2.3 依赖安装

```bash
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/golang-jwt/jwt/v5
go get -u golang.org/x/crypto/bcrypt
```

---

## 3. 项目结构

### 3.1 目录组织

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # 应用入口
├── internal/
│   ├── config/
│   │   └── config.go              # 配置管理
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
│   └── dto/
│       ├── login_request.go       # 登录请求 DTO
│       ├── register_request.go    # 注册请求 DTO
│       └── auth_response.go       # 认证响应 DTO
├── migrations/
│   ├── 000001_create_users_table.up.sql
│   └── 000001_create_users_table.down.sql
├── pkg/
│   ├── jwt/
│   │   └── jwt.go                 # JWT 工具
│   ├── password/
│   │   └── password.go            # 密码加密工具
│   └── response/
│       └── response.go            # 统一响应格式
├── configs/
│   ├── config.yaml                # 配置文件
│   └── config.development.yaml    # 开发环境配置
├── .env.example                   # 环境变量示例
├── .air.toml                      # Air 热重载配置
├── Dockerfile                     # Docker 镜像构建
├── docker-compose.yml             # Docker Compose 配置
├── go.mod
├── go.sum
└── README.md
```

### 3.2 代码分层架构

```
┌─────────────────────────────────────────┐
│         Handler Layer (HTTP)            │  ← 处理 HTTP 请求/响应
├─────────────────────────────────────────┤
│         Service Layer (Business)        │  ← 业务逻辑
├─────────────────────────────────────────┤
│      Repository Layer (Data Access)     │  ← 数据库操作
├─────────────────────────────────────────┤
│         Database (PostgreSQL)           │  ← 数据持久化
└─────────────────────────────────────────┘
```

---

## 4. 数据库设计

### 4.1 User 表

**表名**: `users`

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | UUID | PRIMARY KEY | 用户唯一标识 |
| username | VARCHAR(50) | NOT NULL, UNIQUE | 用户名 |
| email | VARCHAR(255) | NOT NULL, UNIQUE | 邮箱 |
| password_hash | VARCHAR(255) | NOT NULL | 密码哈希（bcrypt） |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 |

### 4.2 索引

```sql
-- 用户名唯一索引
CREATE UNIQUE INDEX idx_users_username ON users(username);

-- 邮箱唯一索引
CREATE UNIQUE INDEX idx_users_email ON users(email);

-- 创建时间索引（用于排序）
CREATE INDEX idx_users_created_at ON users(created_at DESC);
```

### 4.3 数据库迁移

**文件**: `migrations/000001_create_users_table.up.sql`

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- 创建更新时间触发器
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

**文件**: `migrations/000001_create_users_table.down.sql`

```sql
DROP TRIGGER update_users_updated_at ON users;
DROP FUNCTION update_updated_at_column();
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;
```

---

## 5. API 接口设计

### 5.1 接口列表

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/auth/register | 用户注册 | 否 |
| POST | /api/auth/login | 用户登录 | 否 |
| POST | /api/auth/logout | 用户登出 | 是 |
| POST | /api/auth/refresh | 刷新 Token | 否 (使用 Refresh Token) |
| GET | /api/auth/me | 获取当前用户 | 是 |

### 5.2 数据格式

#### 5.2.1 用户注册

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

**错误响应**: `400 Bad Request`

```json
{
  "code": 400,
  "message": "Passwords do not match",
  "data": null
}
```

**错误响应**: `409 Conflict`

```json
{
  "code": 409,
  "message": "User with this email already exists",
  "data": null
}
```

#### 5.2.2 用户登录

**请求**: `POST /api/auth/login`

```json
{
  "email": "john@example.com",
  "password": "password123",
  "rememberMe": false
}
```

**响应**: `200 OK`

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

**错误响应**: `401 Unauthorized`

```json
{
  "code": 401,
  "message": "Invalid email or password",
  "data": null
}
```

#### 5.2.3 用户登出

**请求**: `POST /api/auth/logout`

Headers:
```
Authorization: Bearer {token}
```

**响应**: `200 OK`

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

#### 5.2.4 刷新 Token

**请求**: `POST /api/auth/refresh`

Headers:
```
Authorization: Bearer {refreshToken}
```

**响应**: `200 OK`

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

**错误响应**: `401 Unauthorized`

```json
{
  "code": 401,
  "message": "Invalid or expired refresh token",
  "data": null
}
```

#### 5.2.5 获取当前用户

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
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john@example.com",
    "createdAt": "2025-12-26T10:30:00Z",
    "updatedAt": "2025-12-26T10:30:00Z"
  }
}
```

**错误响应**: `401 Unauthorized`

```json
{
  "code": 401,
  "message": "Missing or invalid authorization token",
  "data": null
}
```

### 5.3 统一响应码

| 响应码 | HTTP 状态 | 说明 |
|--------|----------|------|
| 0 | 200/201 | 成功 |
| 400 | 400 | 请求参数验证失败 |
| 401 | 401 | 未授权/Token 无效/密码错误 |
| 404 | 404 | 资源不存在 |
| 409 | 409 | 资源冲突（用户已存在） |
| 500 | 500 | 服务器内部错误 |

---

## 6. 核心模块实现

### 6.1 配置管理

**文件**: `internal/config/config.go`

```go
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port            string
	Mode            string // debug, release
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret          string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
	Issuer          string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

func Load() (*Config, error) {
	// 加载 .env 文件（开发环境）
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			Mode:            getEnv("SERVER_MODE", "debug"),
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "artisan"),
			Password: getEnv("DB_PASSWORD", "artisan123"),
			DBName:   getEnv("DB_NAME", "artisan_coder"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			AccessDuration:  1 * time.Hour,
			RefreshDuration: 7 * 24 * time.Hour, // 7 days
			Issuer:          "artisan-coder",
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{getEnv("FRONTEND_URL", "http://localhost:5173")},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Origin", "Content-Type", "Authorization"},
		},
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
```

### 6.2 User 模型

**文件**: `internal/models/user.go`

```go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username     string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
	Email        string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM hook
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
```

### 6.3 JWT 工具

**文件**: `pkg/jwt/jwt.go`

```go
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret          []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
	issuer          string
}

func NewManager(secret string, accessDuration, refreshDuration time.Duration, issuer string) *Manager {
	return &Manager{
		secret:          []byte(secret),
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
		issuer:          issuer,
	}
}

// GenerateTokenPair 生成访问令牌和刷新令牌
func (m *Manager) GenerateTokenPair(userID uuid.UUID, email string) (accessToken, refreshToken string, err error) {
	// 生成 Access Token
	accessToken, err = m.generateToken(userID, email, m.accessDuration)
	if err != nil {
		return "", "", err
	}

	// 生成 Refresh Token
	refreshToken, err = m.generateToken(userID, email, m.refreshDuration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (m *Manager) generateToken(userID uuid.UUID, email string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ValidateToken 验证并解析 Token
func (m *Manager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken 刷新 Token
func (m *Manager) RefreshToken(refreshToken string) (accessToken, newRefreshToken string, err error) {
	claims, err := m.ValidateToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	return m.GenerateTokenPair(claims.UserID, claims.Email)
}
```

### 6.4 密码加密工具

**文件**: `pkg/password/password.go`

```go
package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// bcrypt cost factor
	cost = 12
)

// Hash 对密码进行 bcrypt 哈希
func Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(bytes), nil
}

// Verify 验证密码是否匹配
func Verify(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
```

### 6.5 统一响应格式

**文件**: `pkg/response/response.go`

```go
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`              // 响应码，0 表示成功
	Message string      `json:"message"`           // 响应消息
	Data    interface{} `json:"data"`              // 响应数据，成功时返回数据，失败时为 null
}

const (
	CodeSuccess           = 0     // 成功
	CodeBadRequest       = 400   // 请求参数错误
	CodeUnauthorized     = 401   // 未授权
	CodeNotFound         = 404   // 资源不存在
	CodeConflict         = 409   // 资源冲突
	CodeInternalError    = 500   // 服务器内部错误
)

const (
	MessageSuccess           = "success"
	MessageBadRequest       = "Bad request"
	MessageUnauthorized     = "Unauthorized"
	MessageNotFound         = "Not found"
	MessageConflict         = "Conflict"
	MessageInternalError    = "Internal server error"
)

// Success 成功响应 (200)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: MessageSuccess,
		Data:    data,
	})
}

// Created 创建成功响应 (201)
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    CodeSuccess,
		Message: MessageSuccess,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, statusCode int, code int, message string) {
	c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest 400 错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, CodeBadRequest, message)
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = MessageUnauthorized
	}
	Error(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

// NotFound 404 错误
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = MessageNotFound
	}
	Error(c, http.StatusNotFound, CodeNotFound, message)
}

// Conflict 409 错误
func Conflict(c *gin.Context, message string) {
	if message == "" {
		message = MessageConflict
	}
	Error(c, http.StatusConflict, CodeConflict, message)
}

// InternalError 500 错误
func InternalError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, CodeInternalError, MessageInternalError)
}
```

### 6.6 中间件

#### CORS 中间件

**文件**: `internal/middleware/cors.go`

```go
package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 检查是否允许该源
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
```

#### JWT 认证中间件

**文件**: `internal/middleware/auth.go`

```go
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"artisan-coder/pkg/jwt"
	"artisan-coder/pkg/response"
)

const (
	userIDKey = "user_id"
)

func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Missing authorization token")
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// 将用户 ID 存储到上下文
		c.Set(userIDKey, claims.UserID)
		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(userIDKey)
	if !exists {
		return "", false
	}
	return userID.(string), true
}
```

### 6.7 Repository 层

**文件**: `internal/repository/user_repository.go`

```go
package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		// 检查是否是唯一约束冲突
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrUserAlreadyExists
		}
		return result.Error
	}
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *User) error {
	result := r.db.WithContext(ctx).Save(user)
	return result.Error
}
```

### 6.8 Service 层

**文件**: `internal/service/auth_service.go`

```go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"artisan-coder/internal/models"
	"artisan-coder/internal/repository"
	"artisan-coder/pkg/jwt"
	"artisan-coder/pkg/password"
)

type AuthService interface {
	Register(ctx context.Context, username, email, userPassword string) (*models.User, string, string, error)
	Login(ctx context.Context, email, userPassword string) (*models.User, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.User, string, string, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.Manager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.Manager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(ctx context.Context, username, email, userPassword string) (*models.User, string, string, error) {
	// 检查用户是否已存在
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, "", "", repository.ErrUserAlreadyExists
	} else if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, "", "", err
	}

	// 加密密码
	hashedPassword, err := password.Hash(userPassword)
	if err != nil {
		return nil, "", "", err
	}

	// 创建用户
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", "", err
	}

	// 生成 Token
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) Login(ctx context.Context, email, userPassword string) (*models.User, string, string, error) {
	// 查找用户
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, "", "", errors.New("invalid credentials")
		}
		return nil, "", "", err
	}

	// 验证密码
	if !password.Verify(user.PasswordHash, userPassword) {
		return nil, "", "", errors.New("invalid credentials")
	}

	// 生成 Token
	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.User, string, string, error) {
	// 验证并刷新 Token
	accessToken, newRefreshToken, err := s.jwtManager.RefreshToken(refreshToken)
	if err != nil {
		return nil, "", "", errors.New("invalid or expired refresh token")
	}

	// 从 Access Token 中获取用户信息
	claims, err := s.jwtManager.ValidateToken(accessToken)
	if err != nil {
		return nil, "", "", err
	}

	// 获取用户信息
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, newRefreshToken, nil
}

func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
```

### 6.9 Handler 层

**文件**: `internal/handler/auth_handler.go`

```go
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"artisan-coder/internal/service"
	"artisan-coder/pkg/response"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=50"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=7"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"rememberMe"`
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	Token        string        `json:"token"`
	RefreshToken string        `json:"refreshToken"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证密码一致性
	if req.Password != req.ConfirmPassword {
		response.BadRequest(c, "Passwords do not match")
		return
	}

	// 调用服务层
	user, accessToken, refreshToken, err := h.authService.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			response.Conflict(c, "User with this email already exists")
		} else {
			response.InternalError(c)
		}
		return
	}

	response.Created(c, &AuthResponse{
		User:         toUserResponse(user),
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 调用服务层
	user, accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Unauthorized(c, "Invalid email or password")
		return
	}

	response.Success(c, &AuthResponse{
		User:         toUserResponse(user),
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT 无状态，客户端删除 Token 即可
	// 如果需要实现 Token 黑名单，可以在此添加逻辑
	response.Success(c, nil)
}

// RefreshToken 刷新 Token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		response.Unauthorized(c, "Missing refresh token")
		return
	}

	// 移除 "Bearer " 前缀
	if len(refreshToken) > 7 && refreshToken[:7] == "Bearer " {
		refreshToken = refreshToken[7:]
	}

	// 调用服务层
	user, accessToken, newRefreshToken, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired refresh token")
		return
	}

	response.Success(c, &AuthResponse{
		User:         toUserResponse(user),
		Token:        accessToken,
		RefreshToken: newRefreshToken,
	})
}

// GetCurrentUser 获取当前用户
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c)
		return
	}

	response.Success(c, toUserResponse(user))
}

func toUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
```

### 6.10 主程序

**文件**: `cmd/server/main.go`

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"artisan-coder/internal/config"
	"artisan-coder/internal/handler"
	"artisan-coder/internal/middleware"
	"artisan-coder/internal/models"
	"artisan-coder/internal/repository"
	"artisan-coder/internal/service"
	"artisan-coder/pkg/jwt"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化依赖
	jwtManager := jwt.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessDuration,
		cfg.JWT.RefreshDuration,
		cfg.JWT.Issuer,
	)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	router := gin.Default()

	// 全局中间件
	router.Use(middleware.CORS(cfg.CORS.AllowedOrigins))
	router.Use(middleware.Logger())

	// 注册路由
	setupRoutes(router, authHandler, jwtManager)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器
	go func() {
		log.Printf("Server is running on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移（开发环境）
	// 生产环境应使用 golang-migrate
	if os.Getenv("APP_ENV") != "production" {
		if err := db.AutoMigrate(&models.User{}); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func setupRoutes(router *gin.Engine, authHandler *handler.AuthHandler, jwtManager *jwt.Manager) {
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.RefreshToken)

			// 需要认证的路由
			auth.GET("/me", middleware.Auth(jwtManager), authHandler.GetCurrentUser)
		}
	}
}
```

---

## 7. 开发任务分解

### 7.1 任务列表

| 阶段 | 任务 | 优先级 |
|------|------|--------|
| 2.1 | 项目初始化 | P0 |
| 2.2 | 数据库设置 | P0 |
| 2.3 | 核心工具实现 | P0 |
| 2.4 | 数据模型和 Repository | P0 |
| 2.5 | 业务逻辑层 | P0 |
| 2.6 | HTTP 层 | P0 |
| 2.7 | 中间件实现 | P0 |
| 2.8 | 主程序和路由 | P0 |
| 2.9 | 测试 | P1 |
| 2.10 | Docker 部署 | P1 |

### 7.2 详细任务

#### 2.1 项目初始化

**任务**:
- [ ] 创建 `backend/` 目录
- [ ] 初始化 Go module: `go mod init artisan-coder`
- [ ] 创建项目目录结构
- [ ] 创建 `.env.example` 文件

**环境变量示例**: `.env.example`

```bash
# Server
SERVER_PORT=8080
SERVER_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=artisan
DB_PASSWORD=artisan123
DB_NAME=artisan_coder
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key-change-in-production-use-openssl-rand-base64-32

# CORS
FRONTEND_URL=http://localhost:5173

# App
APP_ENV=development
```

---

#### 2.2 数据库设置

**任务**:
- [ ] 安装 PostgreSQL（本地或 Docker）
- [ ] 创建数据库: `artisan_coder`
- [ ] 安装 golang-migrate
- [ ] 创建迁移文件: `migrations/000001_create_users_table.up.sql`
- [ ] 创建迁移回滚文件: `migrations/000001_create_users_table.down.sql`
- [ ] 运行迁移: `migrate -path migrations -database "postgres://..." up`

**Docker PostgreSQL**:

```bash
docker run --name artisan-postgres \
  -e POSTGRES_USER=artisan \
  -e POSTGRES_PASSWORD=artisan123 \
  -e POSTGRES_DB=artisan_coder \
  -p 5432:5432 \
  -d postgres:15-alpine
```

---

#### 2.3 核心工具实现

**任务**:
- [ ] 实现 `pkg/jwt/jwt.go` - JWT 工具
- [ ] 实现 `pkg/password/password.go` - 密码加密
- [ ] 实现 `pkg/response/response.go` - 统一响应
- [ ] 实现 `internal/config/config.go` - 配置管理

**验证**:
- JWT 生成和验证测试
- 密码哈希和验证测试

---

#### 2.4 数据模型和 Repository

**任务**:
- [ ] 实现 `internal/models/user.go` - User 模型
- [ ] 实现 `internal/repository/user_repository.go` - Repository 层
- [ ] 编写 Repository 单元测试

**验证**:
- 创建用户成功
- 通过邮箱查找用户
- 通过 ID 查找用户
- 用户不存在时返回错误

---

#### 2.5 业务逻辑层

**任务**:
- [ ] 实现 `internal/service/auth_service.go` - Service 层
- [ ] 实现注册逻辑（密码加密、用户创建、Token 生成）
- [ ] 实现登录逻辑（密码验证、Token 生成）
- [ ] 实现 Token 刷新逻辑
- [ ] 实现 GetCurrentUser 逻辑
- [ ] 编写 Service 单元测试

**验证**:
- 注册成功返回用户和 Token
- 登录成功返回用户和 Token
- 密码错误时登录失败
- 用户不存在时登录失败
- Token 刷新正常工作

---

#### 2.6 HTTP 层

**任务**:
- [ ] 实现 `internal/dto/` 目录下的 DTO 结构
- [ ] 实现 `internal/handler/auth_handler.go` - Handler 层
- [ ] 实现所有 API 端点（Register, Login, Logout, Refresh, GetCurrentUser）
- [ ] 实现请求参数验证
- [ ] 实现 JSON 响应格式

**验证**:
- 所有 API 端点返回正确的状态码
- 错误响应格式符合规范
- 成功响应包含正确的数据

---

#### 2.7 中间件实现

**任务**:
- [ ] 实现 `internal/middleware/cors.go` - CORS 中间件
- [ ] 实现 `internal/middleware/auth.go` - JWT 认证中间件
- [ ] 实现 `internal/middleware/logger.go` - 日志中间件

**验证**:
- CORS 正确设置响应头
- 未认证用户无法访问受保护路由
- 认证用户可以访问受保护路由

---

#### 2.8 主程序和路由

**任务**:
- [ ] 实现 `cmd/server/main.go` - 主程序入口
- [ ] 配置路由组 (`/api/auth`)
- [ ] 实现优雅关闭
- [ ] 实现数据库连接初始化
- [ ] 实现依赖注入

**验证**:
- 服务器可以正常启动
- 所有路由可以正常访问
- 服务器可以优雅关闭

---

#### 2.9 测试

**任务**:
- [ ] 编写 Repository 单元测试
- [ ] 编写 Service 单元测试
- [ ] 编写 Handler 集成测试
- [ ] 测试错误场景
- [ ] 测试边界情况

**测试覆盖率目标**: > 80%

---

#### 2.10 Docker 部署

**任务**:
- [ ] 创建 `Dockerfile`
- [ ] 创建 `docker-compose.yml`
- [ ] 创建 `.dockerignore`
- [ ] 测试 Docker 构建
- [ ] 测试 Docker Compose 启动

**验证**:
- Docker 镜像可以成功构建
- `docker-compose up` 可以启动所有服务
- 前端可以访问后端 API

---

## 8. 测试计划

### 8.1 单元测试

#### Repository 测试

**文件**: `internal/repository/user_repository_test.go`

```go
package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"artisan-coder/internal/models"
	"artisan-coder/internal/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db.AutoMigrate(&models.User{})
	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		PasswordHash: "hashed_password",
	}

	err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	// 创建测试用户
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		PasswordHash: "hashed_password",
	}
	err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// 查找用户
	found, err := repo.FindByEmail(context.Background(), "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, found.Email)

	// 查找不存在的用户
	_, err = repo.FindByEmail(context.Background(), "nonexistent@example.com")
	assert.Error(t, err)
	assert.Equal(t, repository.ErrUserNotFound, err)
}
```

#### Service 测试

**文件**: `internal/service/auth_service_test.go`

```go
package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"artisan-coder/internal/models"
	"artisan-coder/internal/repository"
	"artisan-coder/internal/service"
	"artisan-coder/pkg/jwt"
)

func setupTestService(t *testing.T) (service.AuthService, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	db.AutoMigrate(&models.User{})

	userRepo := repository.NewUserRepository(db)
	jwtManager := jwt.NewManager("test-secret", time.Hour, 7*24*time.Hour, "test")
	authService := service.NewAuthService(userRepo, jwtManager)

	return authService, db
}

func TestAuthService_Register(t *testing.T) {
	authService, _ := setupTestService(t)

	user, accessToken, refreshToken, err := authService.Register(
		context.Background(),
		"testuser",
		"test@example.com",
		"password123",
	)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestAuthService_Login(t *testing.T) {
	authService, _ := setupTestService(t)

	// 先注册
	_, _, _, err := authService.Register(context.Background(), "testuser", "test@example.com", "password123")
	require.NoError(t, err)

	// 登录 - 正确密码
	user, accessToken, refreshToken, err := authService.Login(context.Background(), "test@example.com", "password123")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// 登录 - 错误密码
	_, _, _, err = authService.Login(context.Background(), "test@example.com", "wrongpassword")
	assert.Error(t, err)
}
```

### 8.2 集成测试

**文件**: `test/integration/auth_test.go`

```go
package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"artisan-coder/cmd/server"
)

func TestAuthAPI_Register(t *testing.T) {
	// 启动测试服务器
	router := setupRouter() // 需要实现

	reqBody := map[string]string{
		"username":        "testuser",
		"email":           "test@example.com",
		"password":        "password123",
		"confirmPassword": "password123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	// 检查新的响应格式
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "success", resp["message"])
	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
	assert.NotEmpty(t, data["refreshToken"])
}

func TestAuthAPI_Login(t *testing.T) {
	router := setupRouter()

	// 先注册
	reqBody := map[string]string{
		"username":        "testuser",
		"email":           "test@example.com",
		"password":        "password123",
		"confirmPassword": "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 登录
	loginBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ = json.Marshal(loginBody)
	req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	// 检查新的响应格式
	assert.Equal(t, float64(0), resp["code"])
	assert.Equal(t, "success", resp["message"])
	assert.NotNil(t, resp["data"])
}
```

### 8.3 API 测试

使用 Postman 或 curl 测试：

#### 注册

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123",
    "confirmPassword": "password123"
  }'
```

#### 登录

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

#### 获取当前用户

```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 9. 部署配置

### 9.1 Dockerfile

**文件**: `backend/Dockerfile`

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
```

### 9.2 docker-compose.yml

**文件**: `docker-compose.yml` (项目根目录)

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: artisan-postgres
    environment:
      POSTGRES_USER: artisan
      POSTGRES_PASSWORD: artisan123
      POSTGRES_DB: artisan_coder
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U artisan"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: artisan-backend
    environment:
      SERVER_PORT: 8080
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: artisan
      DB_PASSWORD: artisan123
      DB_NAME: artisan_coder
      DB_SSLMODE: disable
      JWT_SECRET: ${JWT_SECRET:-your-secret-key-change-in-production}
      FRONTEND_URL: http://localhost:5173
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    command: ["./main"]

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: artisan-frontend
    ports:
      - "5173:5173"
    environment:
      VITE_API_BASE_URL: http://localhost:8080
      VITE_USE_MOCK: "false"

volumes:
  postgres_data:
```

### 9.3 .dockerignore

**文件**: `backend/.dockerignore`

```
*.md
.env
.git
.gitignore
.DS_Store
*.log
vendor/
```

### 9.4 部署命令

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f backend

# 停止服务
docker-compose down

# 停止服务并删除数据
docker-compose down -v
```

---

## 10. 前后端联调

### 10.1 API 响应格式说明

**统一响应格式**：

所有 API 接口都使用统一的响应格式：

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

**字段说明**：
- `code`: 响应码，`0` 表示成功，非 `0` 表示失败
- `message`: 响应消息
- `data`: 响应数据，成功时返回实际数据，失败时为 `null`

**响应码定义**：

| 响应码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 500 | 服务器内部错误 |

### 10.2 前端适配

前端需要修改 API 响应拦截器以适配新的响应格式：

**文件**: `frontend/src/services/api.ts`

```typescript
import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor: add token
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token')
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor: unified error handling
apiClient.interceptors.response.use(
  (response) => {
    // 新的响应格式: { code: 0, message: "success", data: {...} }
    const result = response.data

    // 检查 code 是否为 0
    if (result.code === 0) {
      return result.data // 返回 data 字段
    } else {
      // 错误响应
      return Promise.reject({ code: result.code, message: result.message })
    }
  },
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      // Token 过期，清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error.message)
  }
)

export default apiClient
```

**修改认证服务**：

**文件**: `frontend/src/services/auth.ts`

```typescript
import apiClient from './api'
import type { LoginRequest, RegisterRequest } from '@/types/api'
import type { User } from '@/types/models'

interface AuthResponse {
  user: User
  token: string
  refreshToken?: string
}

export const authService = {
  // Login
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/login', data)
    // apiClient 拦截器已经返回 response.data
    // 保存 token 和用户信息
    localStorage.setItem('token', response.token)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },

  // Register
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/register', data)
    localStorage.setItem('token', response.token)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },

  // Logout
  async logout(): Promise<void> {
    try {
      await apiClient.post('/api/auth/logout')
    } finally {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  },

  // Get current user
  async getCurrentUser(): Promise<User> {
    return await apiClient.get<User>('/api/auth/me')
  },

  // Refresh token
  async refreshToken(): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/refresh')
    localStorage.setItem('token', response.token)
    return response
  },
}
```

### 10.3 前端配置

修改前端环境变量：`frontend/.env.development`

```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_USE_MOCK=false
```

### 10.4 联调步骤

1. **启动后端服务**

```bash
cd backend
go run cmd/server/main.go
```

2. **启动前端服务**

```bash
cd frontend
npm run dev
```

3. **测试注册流程**
   - 打开浏览器访问 http://localhost:5173
   - 点击 "Sign up"
   - 填写注册信息
   - 提交表单
   - 验证是否成功注册并跳转到 Dashboard

4. **测试登录流程**
   - 登出
   - 使用注册的账号登录
   - 验证是否成功登录并跳转到 Dashboard

5. **测试 Token 刷新**
   - 等待 Access Token 过期（1小时）
   - 或手动修改 Token 过期时间
   - 验证是否自动刷新 Token

### 10.5 CORS 配置

确保后端 CORS 配置正确允许前端源：

```go
CORS: CORSConfig{
    AllowedOrigins: []string{"http://localhost:5173"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"Origin", "Content-Type", "Authorization"},
}
```

### 10.6 调试技巧

**查看网络请求**:
- 浏览器 DevTools -> Network 标签
- 查看 API 请求的 Request/Response
- 检查响应格式是否为 `{ code: 0, message: "success", data: {...} }`

**查看 Token**:
- 浏览器 DevTools -> Application -> Local Storage
- 查看 `token` 和 `user` 字段

**常见问题**:

1. **CORS 错误**
   - 检查后端 CORS 配置
   - 确认前端 URL 正确

2. **401 Unauthorized**
   - 检查 Token 是否正确发送
   - 查看请求头中的 `Authorization` 字段

3. **响应格式错误**
   - 确认前端 API 拦截器已正确修改
   - 检查后端响应格式是否符合新规范

4. **密码不匹配**
   - 确认前后端密码字段一致
   - 检查密码验证逻辑

---

## 11. 附录

### 11.1 生成 JWT Secret

```bash
# 生成 32 字节的随机密钥
openssl rand -base64 32
```

### 11.2 数据库迁移工具

**安装 golang-migrate**:

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
mv migrate /usr/local/bin/migrate
```

**运行迁移**:

```bash
# 创建迁移文件
migrate create -ext sql -dir migrations -seq create_users_table

# 执行迁移
migrate -path migrations -database "postgres://artisan:artisan123@localhost:5432/artisan_coder?sslmode=disable" up

# 回滚迁移
migrate -path migrations -database "postgres://artisan:artisan123@localhost:5432/artisan_coder?sslmode=disable" down 1
```

### 11.3 开发工具

**Air - 热重载工具**:

**文件**: `.air.toml`

```toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/server"
bin = "tmp/main"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor"]
delay = 1000
stop_on_error = true
```

**安装**:

```bash
go install github.com/cosmtrek/air@latest
```

**运行**:

```bash
air
```

### 11.4 Swagger 文档

**安装 swag**:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

**添加注释**:

```go
// @title Artisan Coder API
// @version 1.0
// @description User authentication API
// @host localhost:8080
// @BasePath /api
func main() { ... }
```

**生成文档**:

```bash
swag init -g cmd/server/main.go
```

### 11.5 常用命令

```bash
# 格式化代码
go fmt ./...

# 代码检查
go vet ./...

# 运行测试
go test ./... -v

# 测试覆盖率
go test ./... -cover

# 构建
go build -o bin/server ./cmd/server

# 运行
./bin/server
```

### 11.6 错误处理最佳实践

1. **使用自定义错误类型**

```go
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}
```

2. **包装错误**

```go
if err != nil {
    return &AppError{
        Code:    "USER_NOT_FOUND",
        Message: "User not found",
        Err:     err,
    }
}
```

3. **日志记录**

```go
import "log"

log.Printf("Error creating user: %v", err)
```

---

## 验收标准

### 功能验收

- [ ] 用户可以成功注册
- [ ] 用户可以成功登录
- [ ] 登录后获取有效的 Access Token 和 Refresh Token
- [ ] Token 过期后可以使用 Refresh Token 刷新
- [ ] 用户可以登出
- [ ] 可以获取当前用户信息
- [ ] 未授权访问受保护路由返回 401

### API 规范验收

- [ ] 所有 API 响应格式符合规范
- [ ] 错误响应包含正确的错误码和消息
- [ ] CORS 正确配置
- [ ] JWT Token 正确生成和验证

### 数据库验收

- [ ] User 表正确创建
- [ ] 密码使用 bcrypt 加密存储
- [ ] 索引正确创建

### 代码质量验收

- [ ] 代码通过 go vet 检查
- [ ] 代码通过 go fmt 格式化
- [ ] 单元测试覆盖率 > 80%
- [ ] 所有测试通过

### 部署验收

- [ ] Docker 镜像可以成功构建
- [ ] docker-compose up 可以启动所有服务
- [ ] 前端可以成功调用后端 API

---

## 下一步计划

完成 Phase 2 后，进入 **Phase 3: 项目管理功能**

- 项目 CRUD API
- 容器生命周期管理
- 文件浏览器 API
- SSH 终端 WebSocket 连接

---

**文档版本**: v1.0
**最后更新**: 2025-12-26
**参考文档**: docs/requirements.md, docs/tasks/phase-1-web-auth-dashboard.md
