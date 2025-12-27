# Phase 2.5: 前后端对接 - 用户认证系统

> **目标**: 将 Phase 1 前端与 Phase 2 后端完成对接，实现真实的用户认证功能
>
> **状态**: 🔄 进行中
> **预计工时**: 2-3 天

---

## 目录

- [1. 对接概述](#1-对接概述)
- [2. API 接口对接](#2-api-接口对接)
- [3. 前端适配](#3-前端适配)
- [4. 后端适配](#4-后端适配)
- [5. 数据格式对齐](#5-数据格式对齐)
- [6. 环境配置](#6-环境配置)
- [7. 测试验证](#7-测试验证)
- [8. 常见问题](#8-常见问题)
- [9. 部署配置](#9-部署配置)

---

## 1. 对接概述

### 1.1 对接目标

将 Phase 1 的前端 Mock 数据替换为真实的 Phase 2 后端 API，实现：

- ✅ 用户注册
- ✅ 用户登录
- ✅ 用户登出
- ✅ Token 刷新
- ✅ 获取当前用户信息
- ✅ 认证中间件保护路由

### 1.2 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| **前端** | React + TypeScript + Vite | Phase 1 前端 |
| **后端** | Go + Gin + GORM | Phase 2 后端 |
| **数据库** | PostgreSQL | 用户数据存储 |
| **认证** | JWT (HS256) | Access Token + Refresh Token |
| **通信** | HTTP/REST API | JSON 格式 |

### 1.3 后端 API 端点

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/auth/register` | 用户注册 | 否 |
| POST | `/api/auth/login` | 用户登录 | 否 |
| POST | `/api/auth/logout` | 用户登出 | 是 |
| POST | `/api/auth/refresh` | 刷新 Token | 否 (使用 Refresh Token) |
| GET | `/api/auth/me` | 获取当前用户 | 是 |

### 1.4 前端现有实现

Phase 1 前端已完成：

- ✅ 登录页面 UI (`frontend/src/pages/Login.tsx`)
- ✅ 注册页面 UI (`frontend/src/pages/Register.tsx`)
- ✅ Dashboard 页面 (`frontend/src/pages/Dashboard.tsx`)
- ✅ 用户状态管理 (`frontend/src/stores/userStore.ts`)
- ✅ Mock API 服务 (`frontend/src/services/mockApi.ts`)
- ✅ 路由守卫 (`frontend/src/components/auth/ProtectedRoute.tsx`)
- ✅ API 客户端配置 (`frontend/src/services/api.ts`)

---

## 2. API 接口对接

### 2.1 后端统一响应格式

**重要**: 所有后端 API 都使用统一的响应格式：

```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

**字段说明**:
- `code`: 响应码，`0` 表示成功，非 `0` 表示失败
- `message`: 响应消息
- `data`: 响应数据，成功时返回实际数据，失败时为 `null`

**响应码定义**:

| 响应码 | HTTP 状态 | 说明 |
|--------|----------|------|
| 0 | 200/201 | 成功 |
| 400 | 400 | 请求参数错误 |
| 401 | 401 | 未授权/Token 无效/密码错误 |
| 404 | 404 | 资源不存在 |
| 409 | 409 | 资源冲突（用户已存在） |
| 500 | 500 | 服务器内部错误 |

### 2.2 用户注册 API

#### 2.2.1 接口定义

**请求**: `POST /api/auth/register`

**请求头**:
```http
Content-Type: application/json
```

**请求体**:
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "confirmPassword": "password123"
}
```

**字段验证规则**:
- `username`: 必填，3-50 字符，只能包含字母、数字、连字符和下划线
- `email`: 必填，有效的邮箱格式
- `password`: 必填，最少 7 个字符
- `confirmPassword`: 必填，必须与 password 一致

**成功响应**: `201 Created`

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

**错误响应示例**:

1. 密码不匹配 (`400 Bad Request`):
```json
{
  "code": 400,
  "message": "Passwords do not match",
  "data": null
}
```

2. 邮箱已存在 (`409 Conflict`):
```json
{
  "code": 409,
  "message": "User with this email already exists",
  "data": null
}
```

#### 2.2.2 前端对接代码

**文件**: `frontend/src/services/auth.ts`

```typescript
// 注册
async register(data: RegisterRequest): Promise<AuthResponse> {
  const response = await apiClient.post<AuthResponse>('/api/auth/register', data)
  // 保存 token 和用户信息
  localStorage.setItem('token', response.token)
  localStorage.setItem('refreshToken', response.refreshToken)
  localStorage.setItem('user', JSON.stringify(response.user))
  return response
}
```

### 2.3 用户登录 API

#### 2.3.1 接口定义

**请求**: `POST /api/auth/login`

**请求头**:
```http
Content-Type: application/json
```

**请求体**:
```json
{
  "email": "john@example.com",
  "password": "password123",
  "rememberMe": false
}
```

**字段说明**:
- `email`: 必填，有效的邮箱格式
- `password`: 必填
- `rememberMe`: 可选，布尔值

**成功响应**: `200 OK`

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

#### 2.3.2 前端对接代码

**文件**: `frontend/src/services/auth.ts`

```typescript
// 登录
async login(data: LoginRequest): Promise<AuthResponse> {
  const response = await apiClient.post<AuthResponse>('/api/auth/login', data)
  // 保存 token 和用户信息
  localStorage.setItem('token', response.token)
  localStorage.setItem('refreshToken', response.refreshToken)
  localStorage.setItem('user', JSON.stringify(response.user))
  return response
}
```

### 2.4 用户登出 API

#### 2.4.1 接口定义

**请求**: `POST /api/auth/logout`

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {token}
```

**成功响应**: `200 OK`

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

#### 2.4.2 前端对接代码

**文件**: `frontend/src/services/auth.ts`

```typescript
// 登出
async logout(): Promise<void> {
  try {
    await apiClient.post('/api/auth/logout')
  } finally {
    // 无论后端是否成功，都清除本地存储
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
  }
}
```

### 2.5 Token 刷新 API

#### 2.5.1 接口定义

**请求**: `POST /api/auth/refresh`

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {refreshToken}
```

**成功响应**: `200 OK`

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

#### 2.5.2 前端对接代码

**文件**: `frontend/src/services/auth.ts`

```typescript
// 刷新 token
async refreshToken(): Promise<AuthResponse> {
  const response = await apiClient.post<AuthResponse>('/api/auth/refresh')
  // 更新本地存储
  localStorage.setItem('token', response.token)
  localStorage.setItem('refreshToken', response.refreshToken)
  localStorage.setItem('user', JSON.stringify(response.user))
  return response
}
```

### 2.6 获取当前用户 API

#### 2.6.1 接口定义

**请求**: `GET /api/auth/me`

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {token}
```

**成功响应**: `200 OK`

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

#### 2.6.2 前端对接代码

**文件**: `frontend/src/services/auth.ts`

```typescript
// 获取当前用户
async getCurrentUser(): Promise<User> {
  return await apiClient.get<User>('/api/auth/me')
}
```

---

## 3. 前端适配

### 3.1 修改 API 客户端拦截器

**文件**: `frontend/src/services/api.ts`

**现有实现** (Mock 版本):

```typescript
import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@/types/api'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器：添加 token
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

// 响应拦截器：统一错误处理
apiClient.interceptors.response.use(
  (response) => response.data,
  (error: AxiosError<ApiResponse>) => {
    if (error.response?.status === 401) {
      // Token 过期，清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data?.error || error.message)
  }
)

export default apiClient
```

**修改后的实现** (对接真实 API):

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

// 请求拦截器：添加 token
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

// 响应拦截器：适配后端统一响应格式
apiClient.interceptors.response.use(
  (response) => {
    // 后端响应格式: { code: 0, message: "success", data: {...} }
    const result = response.data

    // 检查 code 是否为 0
    if (result.code === 0) {
      return result.data // 返回 data 字段
    } else {
      // 错误响应
      return Promise.reject({
        code: result.code,
        message: result.message,
      })
    }
  },
  (error: AxiosError) => {
    // 处理 HTTP 错误
    if (error.response) {
      const status = error.response.status

      // 401 未授权 - 清除本地存储并跳转到登录页
      if (status === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
        localStorage.removeItem('user')
        window.location.href = '/login'
      }

      // 后端错误响应格式: { code: xxx, message: "error", data: null }
      const errorData = error.response.data as any
      return Promise.reject({
        code: errorData?.code || status,
        message: errorData?.message || error.message,
      })
    }

    // 网络错误或其他错误
    return Promise.reject({
      code: 0,
      message: error.message || 'Network error',
    })
  }
)

export default apiClient
```

**关键修改点**:

1. **响应格式适配**: 后端返回 `{ code, message, data }` 格式，拦截器需要提取 `data` 字段
2. **错误处理增强**: 统一处理 `code !== 0` 的情况
3. **Refresh Token 支持**: 添加 `refreshToken` 的存储和清理
4. **401 处理**: Token 过期时自动清除本地存储并跳转登录页

### 3.2 修改认证服务

**文件**: `frontend/src/services/auth.ts`

**完整实现**:

```typescript
import apiClient from './api'
import type { LoginRequest, RegisterRequest } from '@/types/api'
import type { User } from '@/types/models'

interface AuthResponse {
  user: User
  token: string
  refreshToken: string
}

export const authService = {
  // 登录
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/login', data)
    // 保存 token 和用户信息
    localStorage.setItem('token', response.token)
    localStorage.setItem('refreshToken', response.refreshToken)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },

  // 注册
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/register', data)
    // 保存 token 和用户信息
    localStorage.setItem('token', response.token)
    localStorage.setItem('refreshToken', response.refreshToken)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },

  // 登出
  async logout(): Promise<void> {
    try {
      await apiClient.post('/api/auth/logout')
    } finally {
      // 无论后端是否成功，都清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
      localStorage.removeItem('user')
    }
  },

  // 获取当前用户
  async getCurrentUser(): Promise<User> {
    return await apiClient.get<User>('/api/auth/me')
  },

  // 刷新 token
  async refreshToken(): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/refresh')
    // 更新本地存储
    localStorage.setItem('token', response.token)
    localStorage.setItem('refreshToken', response.refreshToken)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },
}
```

**关键修改点**:

1. **添加 Refresh Token 支持**: 所有需要认证的操作都会保存 `refreshToken`
2. **错误处理**: `logout` 方法使用 `try-finally` 确保本地存储总是被清除
3. **Token 刷新**: `refreshToken` 方法会更新所有本地存储的认证信息

### 3.3 修改用户状态管理

**文件**: `frontend/src/stores/userStore.ts`

**需要修改的部分**:

```typescript
import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from '@/types/models'
import type { LoginRequest, RegisterRequest } from '@/types/api'
import { authService } from '@/services/auth'

interface UserState {
  user: User | null
  token: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isLoading: boolean

  // Actions
  login: (credentials: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => Promise<void>
  refreshToken: () => Promise<void>
  initialize: () => void
}

export const useUserStore = create<UserState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      refreshToken: null,  // 新增
      isAuthenticated: false,
      isLoading: false,

      login: async (credentials: LoginRequest) => {
        set({ isLoading: true })
        try {
          const response = await authService.login(credentials)
          set({
            user: response.user,
            token: response.token,
            refreshToken: response.refreshToken,  // 新增
            isAuthenticated: true,
            isLoading: false,
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      register: async (data: RegisterRequest) => {
        set({ isLoading: true })
        try {
          const response = await authService.register(data)
          set({
            user: response.user,
            token: response.token,
            refreshToken: response.refreshToken,  // 新增
            isAuthenticated: true,
            isLoading: false,
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      logout: async () => {
        set({ isLoading: true })
        try {
          await authService.logout()
          set({
            user: null,
            token: null,
            refreshToken: null,  // 新增
            isAuthenticated: false,
            isLoading: false,
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      refreshToken: async () => {
        const { refreshToken: currentRefreshToken } = get()
        if (!currentRefreshToken) {
          throw new Error('No refresh token available')
        }

        try {
          const response = await authService.refreshToken()
          set({
            user: response.user,
            token: response.token,
            refreshToken: response.refreshToken,
          })
        } catch (error) {
          // Token 刷新失败，清除状态
          set({
            user: null,
            token: null,
            refreshToken: null,
            isAuthenticated: false,
          })
          throw error
        }
      },

      initialize: () => {
        const token = localStorage.getItem('token')
        const refreshToken = localStorage.getItem('refreshToken')  // 新增
        const userStr = localStorage.getItem('user')
        const user = userStr ? JSON.parse(userStr) : null

        if (token && user) {
          set({
            user,
            token,
            refreshToken,  // 新增
            isAuthenticated: true,
          })
        }
      },
    }),
    {
      name: 'user-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        refreshToken: state.refreshToken,  // 新增
      }),
    }
  )
)
```

**关键修改点**:

1. **添加 `refreshToken` 字段**: 状态管理中包含 `refreshToken`
2. **持久化配置**: 使用 `partialize` 持久化 `refreshToken`
3. **Token 刷新逻辑**: `refreshToken` 方法在刷新失败时清除状态

### 3.4 添加 Token 自动刷新拦截器

**文件**: `frontend/src/services/api.ts`

**可选优化**: 在 API 客户端中添加自动 Token 刷新逻辑

```typescript
// 创建一个不带认证拦截器的客户端（用于刷新 token）
const authApiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

let isRefreshing = false
let refreshSubscribers: Array<(token: string) => void> = []

function subscribeTokenRefresh(cb: (token: string) => void) {
  refreshSubscribers.push(cb)
}

function onTokenRefreshed(token: string) {
  refreshSubscribers.forEach((cb) => cb(token))
  refreshSubscribers = []
}

// 响应拦截器：自动刷新 Token
apiClient.interceptors.response.use(
  (response) => {
    const result = response.data
    if (result.code === 0) {
      return result.data
    } else {
      return Promise.reject({
        code: result.code,
        message: result.message,
      })
    }
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    if (error.response?.status === 401 && !originalRequest._retry) {
      // Token 过期，尝试刷新
      if (isRefreshing) {
        // 如果正在刷新，将请求加入队列
        return new Promise((resolve) => {
          subscribeTokenRefresh((token) => {
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${token}`
            }
            resolve(apiClient(originalRequest))
          })
        })
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        const refreshToken = localStorage.getItem('refreshToken')
        if (!refreshToken) {
          throw new Error('No refresh token available')
        }

        // 使用 authApiClient 调用刷新接口（避免递归）
        const response = await authApiClient.post('/api/auth/refresh', null, {
          headers: {
            Authorization: `Bearer ${refreshToken}`,
          },
        })

        const result = response.data as { data: { token: string; refreshToken: string } }
        const { token: newToken, refreshToken: newRefreshToken } = result.data

        // 更新本地存储
        localStorage.setItem('token', newToken)
        localStorage.setItem('refreshToken', newRefreshToken)

        // 通知所有等待的请求
        onTokenRefreshed(newToken)

        // 重试原始请求
        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`
        }
        return apiClient(originalRequest)
      } catch (refreshError) {
        // Token 刷新失败，清除本地存储并跳转登录页
        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
        localStorage.removeItem('user')
        window.location.href = '/login'
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    // 其他错误
    if (error.response) {
      const errorData = error.response.data as any
      return Promise.reject({
        code: errorData?.code || error.response.status,
        message: errorData?.message || error.message,
      })
    }

    return Promise.reject({
      code: 0,
      message: error.message || 'Network error',
    })
  }
)
```

**说明**: 这是一个可选的高级功能，可以自动刷新过期的 Token 并重试原始请求。

### 3.5 修改环境变量

**文件**: `frontend/.env.development`

```bash
# API 基础 URL
VITE_API_BASE_URL=http://localhost:8080

# 使用 Mock API (改为 false)
VITE_USE_MOCK=false
```

**文件**: `frontend/.env.production`

```bash
# API 基础 URL (生产环境)
VITE_API_BASE_URL=https://api.artisancoder.com

# 使用 Mock API (改为 false)
VITE_USE_MOCK=false
```

### 3.6 更新类型定义

**文件**: `frontend/src/types/api.ts`

确保类型定义与后端响应格式一致：

```typescript
// 认证相关
export interface LoginRequest {
  email: string
  password: string
  rememberMe?: boolean
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  confirmPassword: string
}

export interface AuthResponse {
  user: User
  token: string
  refreshToken: string
}

// API 错误响应
export interface ApiError {
  code: number
  message: string
}

// 统一 API 响应（后端格式）
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}
```

---

## 4. 后端适配

### 4.1 CORS 配置检查

确保后端 CORS 配置正确允许前端源。

**文件**: `backend/configs/config.development.yaml`

```yaml
cors:
  allowedOrigins:
    - "http://localhost:5173"
    - "http://localhost:3000"
  allowedMethods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowedHeaders:
    - "Origin"
    - "Content-Type"
    - "Authorization"
```

**文件**: `backend/internal/middleware/cors.go`

确保 CORS 中间件实现正确：

```go
package middleware

import (
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

### 4.2 响应格式验证

确保所有 API 响应都符合统一格式。

**文件**: `backend/pkg/response/response.go`

```go
package response

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
  Code    int         `json:"code"`
  Message string      `json:"message"`
  Data    interface{} `json:"data"`
}

// Success 成功响应 (200)
func Success(c *gin.Context, data interface{}) {
  c.JSON(http.StatusOK, Response{
    Code:    0,
    Message: "success",
    Data:    data,
  })
}

// Created 创建成功响应 (201)
func Created(c *gin.Context, data interface{}) {
  c.JSON(http.StatusCreated, Response{
    Code:    0,
    Message: "success",
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
  Error(c, http.StatusBadRequest, 400, message)
}

// Unauthorized 401 错误
func Unauthorized(c *gin.Context, message string) {
  if message == "" {
    message = "Unauthorized"
  }
  Error(c, http.StatusUnauthorized, 401, message)
}

// NotFound 404 错误
func NotFound(c *gin.Context, message string) {
  if message == "" {
    message = "Not found"
  }
  Error(c, http.StatusNotFound, 404, message)
}

// Conflict 409 错误
func Conflict(c *gin.Context, message string) {
  if message == "" {
    message = "Conflict"
  }
  Error(c, http.StatusConflict, 409, message)
}

// InternalError 500 错误
func InternalError(c *gin.Context) {
  Error(c, http.StatusInternalServerError, 500, "Internal server error")
}
```

### 4.3 JWT 配置检查

确保 JWT 配置正确。

**文件**: `backend/configs/config.development.yaml`

```yaml
jwt:
  secret: "development-secret-key-do-not-use-in-production"
  accessDuration: "1h"      # Access Token 有效期：1 小时
  refreshDuration: "168h"   # Refresh Token 有效期：7 天
  issuer: "artisan-coder"
```

**文件**: `backend/pkg/jwt/jwt.go`

确保 JWT Manager 实现正确：

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
```

### 4.4 认证中间件检查

确保认证中间件正确处理 JWT Token。

**文件**: `backend/internal/middleware/auth.go`

```go
package middleware

import (
  "net/http"
  "strings"
  "github.com/gin-gonic/gin"
  "artisan-coder/pkg/jwt"
  "artisan-coder/pkg/response"
)

const userIDKey = "user_id"

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
    c.Set(userIDKey, claims.UserID.String())
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

---

## 5. 数据格式对齐

### 5.1 用户模型对比

#### 前端用户模型

**文件**: `frontend/src/types/models.ts`

```typescript
export interface User {
  id: string
  username: string
  email: string
  createdAt: string
  updatedAt: string
}
```

#### 后端用户模型

**文件**: `backend/internal/models/user.go`

```go
package models

import (
  "time"
  "github.com/google/uuid"
  "gorm.io/gorm"
)

type User struct {
  ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
  Username     string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
  Email        string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
  PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
  CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
  UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
}
```

**对比结果**: ✅ 完全匹配

### 5.2 日期格式对齐

**后端**: 使用 RFC3339 格式 (ISO 8601)
- 示例: `2025-12-26T10:30:00Z`

**前端**: JavaScript `Date` 对象可以自动解析 RFC3339 格式
- `new Date("2025-12-26T10:30:00Z")` ✅

**结论**: ✅ 无需额外转换

### 5.3 字段命名对比

| 字段 | 前端 | 后端 | 一致性 |
|------|------|------|--------|
| 用户 ID | `id` | `id` | ✅ |
| 用户名 | `username` | `username` | ✅ |
| 邮箱 | `email` | `email` | ✅ |
| 创建时间 | `createdAt` | `createdAt` | ✅ |
| 更新时间 | `updatedAt` | `updatedAt` | ✅ |
| Token | `token` | `token` | ✅ |
| Refresh Token | `refreshToken` | `refreshToken` | ✅ |

**结论**: ✅ 所有字段命名一致

---

## 6. 环境配置

### 6.1 开发环境配置

#### 后端配置

**文件**: `backend/configs/config.development.yaml`

```yaml
server:
  port: "8080"
  mode: "debug"
  readTimeout: 15s
  writeTimeout: 15s
  shutdownTimeout: 10s

database:
  host: "localhost"
  port: "5432"
  user: "artisan"
  password: "artisan123"
  dbName: "artisan_coder"
  sslMode: "disable"

jwt:
  secret: "development-secret-key-do-not-use-in-production"
  accessDuration: "1h"
  refreshDuration: "168h"
  issuer: "artisan-coder"

cors:
  allowedOrigins:
    - "http://localhost:5173"
    - "http://localhost:3000"
  allowedMethods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowedHeaders:
    - "Origin"
    - "Content-Type"
    - "Authorization"
```

**环境变量**: `backend/.env`

```bash
APP_ENV=development
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=artisan
DB_PASSWORD=artisan123
DB_NAME=artisan_coder
DB_SSLMODE=disable
JWT_SECRET=development-secret-key-do-not-use-in-production
FRONTEND_URL=http://localhost:5173
```

#### 前端配置

**文件**: `frontend/.env.development`

```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_USE_MOCK=false
```

### 6.2 启动顺序

**1. 启动 PostgreSQL**

```bash
# 方式 1: 使用 Docker
docker run --name artisan-postgres \
  -e POSTGRES_USER=artisan \
  -e POSTGRES_PASSWORD=artisan123 \
  -e POSTGRES_DB=artisan_coder \
  -p 5432:5432 \
  -d postgres:15-alpine

# 方式 2: 使用本地 PostgreSQL
# 确保 PostgreSQL 服务已启动
```

**2. 运行数据库迁移** (如果需要)

```bash
cd backend
migrate -path migrations -database "postgres://artisan:artisan123@localhost:5432/artisan_coder?sslmode=disable" up
```

**3. 启动后端服务**

```bash
cd backend
go run cmd/server/main.go
```

预期输出:
```
[FX] DEBU: Database connected successfully
[FX] DEBU: Database auto migration completed
[FX] DEBU: Starting HTTP server on port 8080
```

**4. 启动前端服务**

```bash
cd frontend
npm run dev
```

预期输出:
```
  VITE v5.x.x  ready in xxx ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

### 6.3 验证连接

**1. 测试后端健康检查** (如果实现了)

```bash
curl http://localhost:8080/health
```

**2. 测试用户注册**

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "confirmPassword": "password123"
  }'
```

预期响应:
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

**3. 测试前端访问**

在浏览器打开: http://localhost:5173

---

## 7. 测试验证

### 7.1 功能测试清单

#### 7.1.1 用户注册测试

**测试步骤**:

1. 打开浏览器访问 http://localhost:5173
2. 点击 "Sign up" 链接
3. 填写注册信息:
   - Username: `testuser`
   - Email: `test@example.com`
   - Password: `password123`
   - Confirm Password: `password123`
4. 点击 "Sign up" 按钮

**预期结果**:

- ✅ 注册成功
- ✅ 自动登录并跳转到 Dashboard
- ✅ Dashboard 显示用户名 "testuser"
- ✅ LocalStorage 中存储了 `token`, `refreshToken`, `user`

**验证方法**:

打开浏览器 DevTools -> Application -> Local Storage，检查:
- `token`: 存在且不为空
- `refreshToken`: 存在且不为空
- `user`: JSON 对象，包含用户信息

打开浏览器 DevTools -> Network，检查:
- 请求: `POST http://localhost:8080/api/auth/register`
- 状态码: `201 Created`
- 响应格式: `{ code: 0, message: "success", data: {...} }`

#### 7.1.2 用户登录测试

**测试步骤**:

1. 点击 "Logout" 登出
2. 在登录页面填写:
   - Email: `test@example.com`
   - Password: `password123`
3. 点击 "Sign in" 按钮

**预期结果**:

- ✅ 登录成功
- ✅ 跳转到 Dashboard
- ✅ Dashboard 显示正确的用户信息

**验证方法**:

打开浏览器 DevTools -> Network，检查:
- 请求: `POST http://localhost:8080/api/auth/login`
- 状态码: `200 OK`
- 请求体: `{ email: "...", password: "..." }`
- 响应格式: `{ code: 0, message: "success", data: {...} }`

#### 7.1.3 认证路由测试

**测试步骤**:

1. 登录成功后，直接在地址栏访问 http://localhost:5173
2. 打开浏览器 DevTools -> Application -> Local Storage
3. 删除 `token` 字段
4. 刷新页面

**预期结果**:

- ✅ 自动跳转到登录页
- ✅ 显示未登录提示

**验证方法**:

打开浏览器 DevTools -> Network，检查:
- 请求: `GET http://localhost:8080/api/auth/me`
- 状态码: `401 Unauthorized`
- 响应: `{ code: 401, message: "...", data: null }`

#### 7.1.4 Token 刷新测试 (可选)

**测试步骤**:

1. 打开浏览器 DevTools -> Application -> Local Storage
2. 修改 `token` 字段为一个无效值（但保留 `refreshToken`）
3. 执行一个需要认证的操作（如访问 Dashboard）

**预期结果**:

- ✅ 自动使用 Refresh Token 刷新
- ✅ 刷新成功后重试原始请求
- ✅ 用户无感知，体验流畅

**验证方法**:

打开浏览器 DevTools -> Network，检查:
- 第一个请求: 失败，返回 `401 Unauthorized`
- 第二个请求: `POST /api/auth/refresh`，成功
- 第三个请求: 重试原始请求，成功

#### 7.1.5 错误处理测试

**测试场景 1: 密码不匹配**

1. 访问注册页面
2. 填写表单，但 `password` 和 `confirmPassword` 不一致
3. 点击 "Sign up"

**预期结果**:

- ✅ 显示错误提示: "Passwords do not match"
- ✅ 不提交请求

**测试场景 2: 邮箱已存在**

1. 尝试使用已存在的邮箱注册

**预期结果**:

- ✅ 显示错误提示: "User with this email already exists"
- ✅ HTTP 状态码: `409 Conflict`

**测试场景 3: 密码错误**

1. 使用错误的密码登录

**预期结果**:

- ✅ 显示错误提示: "Invalid email or password"
- ✅ HTTP 状态码: `401 Unauthorized`

### 7.2 集成测试脚本

**文件**: `frontend/test/integration/auth.spec.ts`

```typescript
import { test, expect } from '@playwright/test'

test.describe('User Authentication', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:5173')
  })

  test('should register a new user', async ({ page }) => {
    // 导航到注册页面
    await page.click('text=Sign up')

    // 填写注册表单
    const timestamp = Date.now()
    const username = `testuser${timestamp}`
    const email = `test${timestamp}@example.com`

    await page.fill('input[name="username"]', username)
    await page.fill('input[name="email"]', email)
    await page.fill('input[name="password"]', 'password123')
    await page.fill('input[name="confirmPassword"]', 'password123')

    // 提交表单
    await page.click('button[type="submit"]')

    // 验证跳转到 Dashboard
    await expect(page).toHaveURL('http://localhost:5173/')
    await expect(page.locator('text=Dashboard')).toBeVisible()

    // 验证用户信息显示
    await expect(page.locator(`text=${username}`)).toBeVisible()
  })

  test('should login with valid credentials', async ({ page }) => {
    // 导航到登录页面
    await page.click('text=Sign in')

    // 填写登录表单
    await page.fill('input[name="email"]', 'test@example.com')
    await page.fill('input[name="password"]', 'password123')

    // 提交表单
    await page.click('button[type="submit"]')

    // 验证跳转到 Dashboard
    await expect(page).toHaveURL('http://localhost:5173/')
    await expect(page.locator('text=Dashboard')).toBeVisible()
  })

  test('should show error with invalid credentials', async ({ page }) => {
    // 导航到登录页面
    await page.click('text=Sign in')

    // 填写错误的登录信息
    await page.fill('input[name="email"]', 'test@example.com')
    await page.fill('input[name="password"]', 'wrongpassword')

    // 提交表单
    await page.click('button[type="submit"]')

    // 验证错误提示
    await expect(page.locator('text=Invalid email or password')).toBeVisible()
  })

  test('should logout and redirect to login', async ({ page }) => {
    // 先登录
    await page.click('text=Sign in')
    await page.fill('input[name="email"]', 'test@example.com')
    await page.fill('input[name="password"]', 'password123')
    await page.click('button[type="submit"]')

    // 等待跳转到 Dashboard
    await expect(page).toHaveURL('http://localhost:5173/')

    // 登出
    await page.click('button[aria-label="User menu"]')
    await page.click('text=Logout')

    // 验证跳转到登录页
    await expect(page).toHaveURL('http://localhost:5173/login')
  })
})
```

---

## 8. 常见问题

### 8.1 CORS 错误

**问题**:
```
Access to XMLHttpRequest at 'http://localhost:8080/api/auth/login' from origin 'http://localhost:5173' has been blocked by CORS policy
```

**原因**: 后端 CORS 配置不允许前端源

**解决方案**:

1. 检查后端 CORS 配置:
```yaml
# backend/configs/config.development.yaml
cors:
  allowedOrigins:
    - "http://localhost:5173"
```

2. 检查环境变量:
```bash
# backend/.env
FRONTEND_URL=http://localhost:5173
```

3. 重启后端服务

### 8.2 401 Unauthorized

**问题**: 所有需要认证的请求都返回 401

**原因**:
- Token 未发送
- Token 格式错误
- Token 已过期

**解决方案**:

1. 检查 Local Storage 中是否有 `token`:
```javascript
// 浏览器 Console
localStorage.getItem('token')
```

2. 检查请求头:
- 打开 DevTools -> Network
- 点击任意请求
- 查看 Request Headers 中的 `Authorization` 字段
- 应该是: `Bearer eyJhbGciOi...`

3. 检查 Token 是否过期:
```bash
# 解码 JWT Token
echo "eyJhbGciOi..." | jwt decode  # 使用 jwt-cli 工具
# 或使用在线工具: https://jwt.io/
```

### 8.3 响应格式错误

**问题**: 前端无法正确解析后端响应

**原因**: 响应格式不匹配

**解决方案**:

1. 检查后端响应:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' | jq
```

预期输出:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user": {...},
    "token": "...",
    "refreshToken": "..."
  }
}
```

2. 检查前端拦截器:
```typescript
// frontend/src/services/api.ts
apiClient.interceptors.response.use(
  (response) => {
    const result = response.data
    if (result.code === 0) {
      return result.data  // ← 应该返回 result.data
    }
    // ...
  }
)
```

### 8.4 网络错误

**问题**: `Network Error` 或 `ERR_CONNECTION_REFUSED`

**原因**: 后端服务未启动或端口错误

**解决方案**:

1. 检查后端是否启动:
```bash
curl http://localhost:8080/health
# 或
lsof -i :8080
```

2. 检查前端环境变量:
```bash
# frontend/.env.development
VITE_API_BASE_URL=http://localhost:8080  # ← 确保端口正确
```

3. 重启前端服务:
```bash
cd frontend
npm run dev
```

### 8.5 数据库连接错误

**问题**: 后端启动时报错 `failed to connect database`

**原因**: PostgreSQL 未启动或配置错误

**解决方案**:

1. 检查 PostgreSQL 是否启动:
```bash
docker ps | grep postgres
# 或
psql -U artisan -d artisan_coder -h localhost
```

2. 检查数据库配置:
```bash
# backend/.env
DB_HOST=localhost
DB_PORT=5432
DB_USER=artisan
DB_PASSWORD=artisan123
DB_NAME=artisan_coder
```

3. 测试数据库连接:
```bash
psql postgresql://artisan:artisan123@localhost:5432/artisan_coder
```

---

## 9. 部署配置

### 9.1 Docker Compose 部署

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
      APP_ENV: production
      SERVER_PORT: 8080
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: artisan
      DB_PASSWORD: artisan123
      DB_NAME: artisan_coder
      DB_SSLMODE: disable
      JWT_SECRET: ${JWT_SECRET}
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

**启动命令**:

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

### 9.2 生产环境配置

#### 后端生产配置

**文件**: `backend/configs/config.production.yaml`

```yaml
server:
  port: "8080"
  mode: "release"
  readTimeout: 15s
  writeTimeout: 15s
  shutdownTimeout: 30s

database:
  host: ""  # 从环境变量读取
  port: "5432"
  user: ""  # 从环境变量读取
  password: ""  # 从环境变量读取
  dbName: ""  # 从环境变量读取
  sslMode: "require"

jwt:
  secret: ""  # 必须从环境变量设置
  accessDuration: "1h"
  refreshDuration: "168h"
  issuer: "artisan-coder"

cors:
  allowedOrigins:
    - "https://artisancoder.com"
    - "https://www.artisancoder.com"
  allowedMethods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowedHeaders:
    - "Origin"
    - "Content-Type"
    - "Authorization"
```

**环境变量**: `backend/.env.production`

```bash
APP_ENV=production
SERVER_PORT=8080
DB_HOST=your-production-db-host
DB_PORT=5432
DB_USER=artisan
DB_PASSWORD=strong-production-password
DB_NAME=artisan_coder
DB_SSLMODE=require
JWT_SECRET=$(openssl rand -base64 32)  # 生成强密钥
FRONTEND_URL=https://artisancoder.com
```

#### 前端生产配置

**文件**: `frontend/.env.production`

```bash
VITE_API_BASE_URL=https://api.artisancoder.com
VITE_USE_MOCK=false
```

### 9.3 安全检查清单

- ✅ JWT Secret 使用强随机值（至少 32 字节）
- ✅ 数据库密码使用强密码
- ✅ 生产环境启用 HTTPS
- ✅ CORS 仅允许可信源
- ✅ 敏感配置通过环境变量传递
- ✅ 数据库连接使用 SSL (`sslmode: require`)
- ✅ Access Token 有效期适中（1 小时）
- ✅ Refresh Token 有效期合理（7 天）
- ✅ 密码使用 bcrypt 加密（cost factor >= 12）

---

## 10. 验收标准

### 10.1 功能验收

- [ ] 用户可以成功注册新账号
- [ ] 用户可以成功登录
- [ ] 登录后跳转到 Dashboard
- [ ] Dashboard 显示正确的用户信息
- [ ] 用户可以登出
- [ ] 未登录访问受保护页面重定向到登录页
- [ ] Token 过期后自动刷新（可选）
- [ ] 错误提示友好明确

### 10.2 API 验收

- [ ] 所有 API 响应格式符合规范
- [ ] 错误响应包含正确的错误码和消息
- [ ] CORS 正确配置
- [ ] JWT Token 正确生成和验证
- [ ] 认证中间件正常工作

### 10.3 代码质量验收

- [ ] 前端 TypeScript 无类型错误
- [ ] 前端 ESLint 无警告
- [ ] 后端 `go vet` 无警告
- [ ] 后端 `go fmt` 格式化正确
- [ ] 所有测试通过

### 10.4 部署验收

- [ ] Docker 镜像可以成功构建
- [ ] `docker-compose up` 可以启动所有服务
- [ ] 前端可以成功调用后端 API
- [ ] 生产环境配置正确

---

## 11. 附录

### 11.1 测试用户数据

用于测试的示例用户：

```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

### 11.2 API 测试脚本

**文件**: `test-api.sh`

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== Testing User Registration ==="
curl -X POST $BASE_URL/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "confirmPassword": "password123"
  }' | jq

echo -e "\n=== Testing User Login ==="
LOGIN_RESPONSE=$(curl -X POST $BASE_URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')

echo "$LOGIN_RESPONSE" | jq

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token')
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.refreshToken')

echo -e "\n=== Testing Get Current User ==="
curl -X GET $BASE_URL/api/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq

echo -e "\n=== Testing Token Refresh ==="
curl -X POST $BASE_URL/api/auth/refresh \
  -H "Authorization: Bearer $REFRESH_TOKEN" | jq

echo -e "\n=== Testing Logout ==="
curl -X POST $BASE_URL/api/auth/logout \
  -H "Authorization: Bearer $TOKEN" | jq
```

使用方法:
```bash
chmod +x test-api.sh
./test-api.sh
```

### 11.3 调试技巧

#### 11.3.1 查看 JWT Token 内容

在线工具: https://jwt.io/

或使用命令行:
```bash
# 安装 jwt-cli
brew install jwt-cli  # macOS
# 或
cargo install jwt-cli  # Rust

# 解码 Token
echo "eyJhbGciOi..." | jwt decode
```

#### 11.3.2 监控 API 请求

**前端**: React Query DevTools

```typescript
// frontend/src/main.tsx
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
    <ReactQueryDevtools initialIsOpen={false} />
  </React.StrictMode>
)
```

**后端**: Gin Logger 中间件已启用

```go
// backend/internal/middleware/logger.go
func Logger() gin.HandlerFunc {
  return func(c *gin.Context) {
    start := time.Now()
    path := c.Request.URL.Path
    query := c.Request.URL.RawQuery

    c.Next()

    end := time.Now()
    latency := end.Sub(start)

    log.Printf("[%s] %s %s | Status: %d | Latency: %v | IP: %s",
      c.Request.Method,
      path,
      query,
      c.Writer.Status(),
      latency,
      c.ClientIP(),
    )
  }
}
```

#### 11.3.3 数据库查询日志

**开发环境**: GORM 可以打印 SQL 查询

```go
// backend/internal/database/database.go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
  Logger: logger.Default.LogMode(logger.Info),  // 打印所有 SQL
})
```

---

## 12. 下一步计划

完成 Phase 2.5 后，进入 **Phase 3: 项目管理功能**

- 项目 CRUD API
- 容器生命周期管理
- 文件浏览器 API
- SSH 终端 WebSocket 连接

---

**文档版本**: v1.0
**创建日期**: 2025-12-27
**依赖文档**:
- docs/requirements.md
- docs/tasks/phase-1-web-auth-dashboard.md
- docs/tasks/phase-2-api-auth-backend.md
