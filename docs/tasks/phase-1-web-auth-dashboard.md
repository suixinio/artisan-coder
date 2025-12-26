# Phase 1: Web 端登录与 Dashboard

> **目标**: 实现用户登录、注册功能，以及登录后的 Dashboard 首页，展示项目列表

**状态**: 🔄 进行中
**预计工时**: 5-7 天

---

## 任务概览

| 阶段 | 任务数 | 状态 |
|------|--------|------|
| 1.1 项目基础 | 5 | ⏳ 待开始 |
| 1.2 基础组件 | 6 | ⏳ 待开始 |
| 1.3 核心服务 | 4 | ⏳ 待开始 |
| 1.4 状态管理 | 2 | ⏳ 待开始 |
| 1.5 认证页面 | 3 | ⏳ 待开始 |
| 1.6 Dashboard | 5 | ⏳ 待开始 |
| 1.7 集成测试 | 3 | ⏳ 待开始 |

---

## 1.1 项目基础搭建

### 1.1.1 创建 Vite + React + TypeScript 项目

**命令**:
```bash
cd /workspaces/artisan-coder
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install
```

**验证**:
```bash
npm run dev
# 访问 http://localhost:5173 确认页面正常
```

**输出**:
- `frontend/` 目录
- 可运行的 React 应用

---

### 1.1.2 安装和配置依赖包

**安装核心依赖**:
```bash
cd frontend

# 路由和状态管理
npm install react-router-dom zustand

# HTTP 客户���
npm install axios

# UI 相关
npm install tailwindcss postcss autoprefixer class-variance-authority clsx tailwind-merge

# Markdown 和代码高亮
npm install react-markdown react-syntax-highlighter

# 图标
npm install lucide-react

# 日期处理
npm install date-fns

# 表单处理
npm install react-hook-form zod @hookform/resolvers
```

**安装开发依赖**:
```bash
npm install -D @types/react-syntax-highlighter prettier eslint-config-prettier
```

---

### 1.1.3 配置 Tailwind CSS

**文件**: `frontend/tailwind.config.js`
```js
/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ['class'],
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        // ... more colors
      },
    },
  },
  plugins: [],
}
```

**文件**: `frontend/src/index.css`
```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --background: 0 0% 100%;
    --foreground: 222.2 84% 4.9%;
    --primary: 221.2 83.2% 53.3%;
    --primary-foreground: 210 40% 98%;
    /* ... more variables */
  }
}
```

---

### 1.1.4 初始化 shadcn/ui

**命令**:
```bash
npx shadcn@latest init
```

**配置选项**:
- TypeScript: Yes
- Style: Default
- Base color: Slate
- CSS variables: Yes
- React Server Components: No
- Import alias: @/*

**文件**: `frontend/components.json`
```json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "default",
  "rsc": false,
  "tsx": true,
  "tailwind": {
    "config": "tailwind.config.js",
    "css": "src/index.css",
    "baseColor": "slate",
    "cssVariables": true
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/utils"
  }
}
```

---

### 1.1.5 配置开发工具

**文件**: `frontend/vite.config.ts`
```ts
import path from 'path'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
```

**文件**: `frontend/tsconfig.json`
```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

**文件**: `frontend/.prettierrc`
```json
{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
```

---

## 1.2 基础 UI 组件

### 1.2.1 安装 shadcn/ui 基础组件

```bash
npx shadcn@latest add button input card label dialog toast dropdown-menu avatar badge separator
```

**输出组件**:
- `src/components/ui/button.tsx`
- `src/components/ui/input.tsx`
- `src/components/ui/card.tsx`
- `src/components/ui/label.tsx`
- `src/components/ui/dialog.tsx`
- `src/components/ui/toast.tsx`
- `src/components/ui/use-toast.ts`
- `src/components/ui/dropdown-menu.tsx`
- `src/components/ui/avatar.tsx`
- `src/components/ui/badge.tsx`
- `src/components/ui/separator.tsx`

---

### 1.2.2 创建工具函数

**文件**: `frontend/src/utils/cn.ts`
```ts
import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
```

---

### 1.2.3 创建 AppShell 布局组件

**文件**: `frontend/src/components/layout/AppShell.tsx`

**功能**:
- 响应式布局容器
- 移动端侧边栏抽屉式
- 桌面端固定侧边栏
- Header 固定顶部
- Main 内容区域可滚动

**接口**:
```tsx
interface AppShellProps {
  children: React.ReactNode
}
```

---

### 1.2.4 创建 Header 组件

**文件**: `frontend/src/components/layout/Header.tsx`

**功能**:
- 显示 Logo 和项目名称
- 移动端汉堡菜单
- 用户菜单 (下拉菜单)
- 登出按钮

**接口**:
```tsx
interface HeaderProps {
  onMenuClick?: () => void
  showMenuButton?: boolean
}
```

---

### 1.2.5 创建 Sidebar 组件

**文件**: `frontend/src/components/layout/Sidebar.tsx`

**功能**:
- 导航菜单 (Dashboard、设置等)
- 可折叠
- 当前路由高亮

**接口**:
```tsx
interface SidebarProps {
  isOpen?: boolean
  onClose?: () => void
}
```

---

### 1.2.6 配置 Toast Provider

**文件**: `frontend/src/components/layout/ToastProvider.tsx`

```tsx
import { Toaster } from '@/components/ui/toast'

export function ToastProvider() {
  return <Toaster />
}
```

**在 main.tsx 中使用**:
```tsx
import { ToastProvider } from '@/components/layout/ToastProvider'

<React.StrictMode>
  <ToastProvider />
  <App />
</React.StrictMode>
```

---

## 1.3 核心服务层

### 1.3.1 定义类型

**文件**: `frontend/src/types/models.ts`

```typescript
// 用户
export interface User {
  id: string
  username: string
  email: string
  createdAt: string
  updatedAt: string
}

// 项目
export interface Project {
  id: string
  userId: string
  name: string
  description?: string
  projectType: 'empty' | 'repo'
  gitRepo?: string
  gitBranch?: string
  containerId?: string
  containerName?: string
  status: 'stopped' | 'running' | 'error'
  sshPort?: number
  createdAt: string
  updatedAt: string
}

// 项目统计
export interface ProjectStats {
  total: number
  running: number
  stopped: number
  error: number
}
```

**文件**: `frontend/src/types/api.ts`

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
  refreshToken?: string
}

// API 响应包装
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: {
    code: string
    message: string
  }
}

// 分页响应
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
}
```

---

### 1.3.2 创建 API 客户端配置

**文件**: `frontend/src/services/api.ts`

```typescript
import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios'
import { ApiResponse } from '@/types/api'

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

---

### 1.3.3 创建认证服务

**文件**: `frontend/src/services/auth.ts`

```typescript
import apiClient from './api'
import { LoginRequest, RegisterRequest, AuthResponse, User } from '@/types/api'

export const authService = {
  // 登录
  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/login', data)
    // 保存 token 和用户信息
    localStorage.setItem('token', response.token)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },

  // 注册
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/register', data)
    localStorage.setItem('token', response.token)
    localStorage.setItem('user', JSON.stringify(response.user))
    return response
  },

  // 登出
  async logout(): Promise<void> {
    try {
      await apiClient.post('/api/auth/logout')
    } finally {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  },

  // 获取当前用户
  async getCurrentUser(): Promise<User> {
    return apiClient.get<User>('/api/auth/me')
  },

  // 刷新 token
  async refreshToken(): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/api/auth/refresh')
    localStorage.setItem('token', response.token)
    return response
  },
}
```

---

### 1.3.4 创建 Mock API 服务

**文件**: `frontend/src/services/mockApi.ts`

```typescript
import { User, Project, ProjectStats } from '@/types/models'

// Mock 用户数据
const mockUsers: User[] = []

// Mock 项目数据
const mockProjects: Project[] = [
  {
    id: '1',
    userId: 'user-1',
    name: 'Example Project',
    description: 'An example project',
    projectType: 'repo',
    gitRepo: 'https://github.com/example/repo',
    gitBranch: 'main',
    status: 'running',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
]

// Mock 延迟
const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

export const mockAuthService = {
  async login(email: string, password: string) {
    await delay(500)

    const user: User = {
      id: 'user-1',
      username: email.split('@')[0],
      email,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }

    const token = btoa(JSON.stringify({ user, exp: Date.now() + 3600000 }))

    return { user, token }
  },

  async register(username: string, email: string, password: string) {
    await delay(500)
    return this.login(email, password)
  },

  async getCurrentUser() {
    await delay(200)
    const token = localStorage.getItem('token')
    if (!token) throw new Error('No token found')

    const parsed = JSON.parse(atob(token))
    return parsed.user
  },
}

export const mockProjectService = {
  async getProjects(): Promise<Project[]> {
    await delay(300)
    return mockProjects
  },

  async getProjectStats(): Promise<ProjectStats> {
    await delay(200)
    return {
      total: mockProjects.length,
      running: mockProjects.filter((p) => p.status === 'running').length,
      stopped: mockProjects.filter((p) => p.status === 'stopped').length,
      error: mockProjects.filter((p) => p.status === 'error').length,
    }
  },
}
```

---

## 1.4 状态管理

### 1.4.1 创建 userStore

**文件**: `frontend/src/stores/userStore.ts`

```typescript
import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { User } from '@/types/models'
import { LoginRequest, RegisterRequest } from '@/types/api'
import { authService, mockAuthService } from '@/services/auth'

// 使用 Mock 还是真实 API
const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true'

interface UserState {
  user: User | null
  token: string | null
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
      isAuthenticated: false,
      isLoading: false,

      login: async (credentials: LoginRequest) => {
        set({ isLoading: true })
        try {
          const response = USE_MOCK
            ? await mockAuthService.login(credentials.email, credentials.password)
            : await authService.login(credentials)

          set({
            user: response.user,
            token: response.token,
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
          const response = USE_MOCK
            ? await mockAuthService.register(data.username, data.email, data.password)
            : await authService.register(data)

          set({
            user: response.user,
            token: response.token,
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
          if (USE_MOCK) {
            localStorage.removeItem('token')
            localStorage.removeItem('user')
          } else {
            await authService.logout()
          }
          set({
            user: null,
            token: null,
            isAuthenticated: false,
            isLoading: false,
          })
        } catch (error) {
          set({ isLoading: false })
          throw error
        }
      },

      refreshToken: async () => {
        const { token } = get()
        if (!token) return

        try {
          const response = await authService.refreshToken()
          set({
            user: response.user,
            token: response.token,
          })
        } catch (error) {
          set({ user: null, token: null, isAuthenticated: false })
          throw error
        }
      },

      initialize: () => {
        const token = localStorage.getItem('token')
        const userStr = localStorage.getItem('user')
        const user = userStr ? JSON.parse(userStr) : null

        if (token && user) {
          set({
            user,
            token,
            isAuthenticated: true,
          })
        }
      },
    }),
    {
      name: 'user-storage',
      partialize: (state) => ({ user: state.user, token: state.token }),
    }
  )
)
```

---

### 1.4.2 创建 projectStore

**文件**: `frontend/src/stores/projectStore.ts`

```typescript
import { create } from 'zustand'
import { Project, ProjectStats } from '@/types/models'
import { mockProjectService } from '@/services/mockApi'

interface ProjectState {
  projects: Project[]
  stats: ProjectStats | null
  isLoading: boolean
  error: string | null

  fetchProjects: () => Promise<void>
  fetchStats: () => Promise<void>
}

export const useProjectStore = create<ProjectState>((set) => ({
  projects: [],
  stats: null,
  isLoading: false,
  error: null,

  fetchProjects: async () => {
    set({ isLoading: true, error: null })
    try {
      const projects = await mockProjectService.getProjects()
      set({ projects, isLoading: false })
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false })
    }
  },

  fetchStats: async () => {
    set({ isLoading: true, error: null })
    try {
      const stats = await mockProjectService.getProjectStats()
      set({ stats, isLoading: false })
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false })
    }
  },
}))
```

---

## 1.5 认证页面

### 1.5.1 创建登录页面

**文件**: `frontend/src/pages/Login.tsx`

**功能**:
- 邮箱/密码输入表单
- 记住我 复选框
- 登录按钮 (loading 状态)
- 表单验证
- 错误提示
- 链接到注册页
- 成功后跳转到 dashboard

**关键代码**:
```tsx
const { login, isLoading } = useUserStore()
const navigate = useNavigate()
const { toast } = useToast()

const form = useForm<LoginRequest>({
  resolver: zodResolver(loginSchema),
  defaultValues: {
    email: '',
    password: '',
  },
})

const onSubmit = async (data: LoginRequest) => {
  try {
    await login(data)
    toast({ title: '登录成功' })
    navigate('/')
  } catch (error) {
    toast({
      title: '登录失败',
      description: (error as Error).message,
      variant: 'destructive',
    })
  }
}
```

---

### 1.5.2 创建注册页面

**文件**: `frontend/src/pages/Register.tsx`

**功能**:
- 用户名、邮箱、密码、确认密码输入
- 密码强度提示
- 表单验证 (两次密码一致)
- 注册按钮 (loading 状态)
- 错误提示
- 链接到登录页
- 成功后自动登录并跳转

---

### 1.5.3 创建路由守卫组件

**文件**: `frontend/src/components/auth/ProtectedRoute.tsx`

```tsx
import { Navigate } from 'react-router-dom'
import { useUserStore } from '@/stores/userStore'

interface ProtectedRouteProps {
  children: React.ReactNode
}

export function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated } = useUserStore()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}
```

**文件**: `frontend/src/components/auth/PublicRoute.tsx`

```tsx
import { Navigate } from 'react-router-dom'
import { useUserStore } from '@/stores/userStore'

interface PublicRouteProps {
  children: React.ReactNode
}

export function PublicRoute({ children }: PublicRouteProps) {
  const { isAuthenticated } = useUserStore()

  if (isAuthenticated) {
    return <Navigate to="/" replace />
  }

  return <>{children}</>
}
```

---

## 1.6 Dashboard 页面

### 1.6.1 创建 Dashboard 主页面

**文件**: `frontend/src/pages/Dashboard.tsx`

**布局**:
```
┌──────────────────────────────────────────────────┐
│  Header: Artisan Coder | 用户菜单 | 登出         │
├──────────┬───────────────────────────────────────┤
│          │  欢迎回来，{username}                  │
│ Sidebar  │  ─────────────────────────────────────│
│          │  统计卡片：                            │
│ - Dashboard│  [总项目: 5] [运行中: 3] [已停止: 2] │
│ - 项目    │                                        │
│ - 设置    │  我的项目                              │
│          │  ┌─────────┐ ┌─────────┐ ┌─────────┐│
│          │  │项目 1   │ │项目 2   │ │项目 3   ││
│          │  │Running  │ │Stopped  │ │Running  ││
│          │  └─────────┘ └─────────┘ └─────────┘│
│          │                                        │
│          │  [+ 创建新项目]                        │
└──────────┴───────────────────────────────────────┘
```

---

### 1.6.2 创建 StatsCard 组件

**文件**: `frontend/src/components/dashboard/StatsCard.tsx`

```tsx
interface StatsCardProps {
  title: string
  value: number
  icon: React.ReactNode
  description?: string
}
```

---

### 1.6.3 创建 ProjectCard 组件

**文件**: `frontend/src/components/project/ProjectCard.tsx`

**功能**:
- 显示项目名称、描述
- 状态徽章
- Git 仓库信息
- 创建时间
- 操作菜单 (启动/停止/打开/删除)

**接口**:
```tsx
interface ProjectCardProps {
  project: Project
  onStart?: () => void
  onStop?: () => void
  onOpen?: () => void
  onDelete?: () => void
}
```

---

### 1.6.4 创建 EmptyState 组件

**文件**: `frontend/src/components/dashboard/EmptyState.tsx`

**功能**:
- 空状态插图或图标
- 提示文字 ("还没有项目")
- "创建第一个项目" 按钮

---

### 1.6.5 配置路由

**文件**: `frontend/src/router.tsx`

```tsx
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { Login, Register, Dashboard } from '@/pages'
import { ProtectedRoute, PublicRoute } from '@/components/auth'
import { AppShell } from '@/components/layout'

const router = createBrowserRouter([
  {
    path: '/login',
    element: (
      <PublicRoute>
        <Login />
      </PublicRoute>
    ),
  },
  {
    path: '/register',
    element: (
      <PublicRoute>
        <Register />
      </PublicRoute>
    ),
  },
  {
    path: '/',
    element: (
      <ProtectedRoute>
        <AppShell>
          <Dashboard />
        </AppShell>
      </ProtectedRoute>
    ),
  },
])

export function Router() {
  return <RouterProvider router={router} />
}
```

---

## 1.7 集成测试

### 1.7.1 端到端流程测试

**测试场景**:
1. 访问应用，未登录自动跳转到登录页
2. 点击注册链接，填写注册信息
3. 注册成功后自动登录并跳转到 dashboard
4. Dashboard 显示用户信息和项目列表
5. 点击登出，返回登录页

### 1.7.2 响应式布局测试

**测试尺寸**:
- Mobile: 375px
- Tablet: 768px
- Desktop: 1024px+
- 4K: 2560px

**检查项**:
- 侧边栏在移动端是否正确隐藏/显示
- Header 是否响应式
- 项目卡片网格是否自适应

### 1.7.3 Mock 数据测试

**测试内容**:
- Mock 登录是否正常工作
- Mock 项目数据是否正确加载
- 错误状态是否正确显示

---

## 验收标准

### 功能验收

- [ ] 用户可以注册新账号
- [ ] 用户可以登录
- [ ] 登录后跳转到 Dashboard
- [ ] Dashboard 显示用户信息
- [ ] Dashboard 显示项目列表 (Mock 数据)
- [ ] Dashboard 显示项目统计
- [ ] 用户可以登出
- [ ] 未登录访问受保护页面重定向到登录页

### UI/UX 验收

- [ ] 界面美观，符合 shadcn/ui 设计风格
- [ ] 响应式布局在移动端、平板、桌面正常显示
- [ ] 加载状态有 spinner 或骨架屏
- [ ] 错误提示友好明确
- [ ] 成功操作有 toast 提示

### 代码质量

- [ ] TypeScript 无类型错误
- [ ] ESLint 无警告
- [ ] 组件职责单一，可复用
- [ ] 有适当的错误处理

---

## 环境变量配置

**文件**: `frontend/.env.development`

```bash
VITE_API_BASE_URL=http://localhost:8080
VITE_USE_MOCK=true
```

**文件**: `frontend/.env.production`

```bash
VITE_API_BASE_URL=https://api.artisancoder.com
VITE_USE_MOCK=false
```

---

## 下一步计划

完成 Phase 1 后，进入 **Phase 2: 项目管理功能**

- 创建项目对话框
- 项目详情页
- 容器启动/停止功能
- 文件浏览器
- SSH 终端
