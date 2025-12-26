# Phase 1: Web 端登录与 Dashboard

> **目标**: 实现用户登录、注册功能，以及登录后的 Dashboard 首页，展示项目列表
>
> **设计参考**: shadcn-admin (New York 风格 + Slate 配色)

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

---

# 附录 A: 设计规范

## A.1 主题配色方案

### 颜色系统
- **色彩空间**: OKLCH (感知均匀，比 HSL 更现代)
- **风格**: New York (shadcn/ui 变体)
- **基础色**: Slate (石板灰/紫灰调)
- **圆角**: 0.625rem (10px)

### 亮色模式色彩变量

**文件**: `frontend/src/styles/theme.css`

```css
:root {
  /* 基础半径 */
  --radius: 0.625rem;

  /* 背景和前景 */
  --background: oklch(1 0 0);                    /* 纯白 #FFFFFF */
  --foreground: oklch(0.129 0.042 264.695);      /* 深紫灰 #1a1a2e */

  /* 卡片和弹出层 */
  --card: oklch(1 0 0);                          /* 纯白 */
  --card-foreground: oklch(0.129 0.042 264.695); /* 深紫灰 */
  --popover: oklch(1 0 0);                       /* 纯白 */
  --popover-foreground: oklch(0.129 0.042 264.695);

  /* 主色调 - 深紫蓝 */
  --primary: oklch(0.208 0.042 265.755);         /* #2d2b55 */
  --primary-foreground: oklch(0.984 0.003 247.858); /* 浅灰白 */

  /* 次要色 */
  --secondary: oklch(0.968 0.007 247.896);       /* 浅灰紫 */
  --secondary-foreground: oklch(0.208 0.042 265.755);

  /* 柔和色 */
  --muted: oklch(0.968 0.007 247.896);          /* 浅灰紫 */
  --muted-foreground: oklch(0.554 0.046 257.417); /* 中灰紫 */

  /* 强调色 */
  --accent: oklch(0.968 0.007 247.896);         /* 浅灰紫 */
  --accent-foreground: oklch(0.208 0.042 265.755);

  /* 危险色 */
  --destructive: oklch(0.577 0.245 27.325);     /* 红橙色 #ef4444 */

  /* 边框和输入 */
  --border: oklch(0.929 0.013 255.508);         /* 浅灰边框 */
  --input: oklch(0.929 0.013 255.508);          /* 输入框边框 */
  --ring: oklch(0.704 0.04 256.788);            /* 焦点环 */

  /* 图表颜色 */
  --chart-1: oklch(0.646 0.222 41.116);         /* 橙色 */
  --chart-2: oklch(0.6 0.118 184.704);          /* 青色 */
  --chart-3: oklch(0.398 0.07 227.392);         /* 蓝色 */
  --chart-4: oklch(0.828 0.189 84.429);         /* 黄色 */
  --chart-5: oklch(0.769 0.188 70.08);          /* 绿色 */

  /* 侧边栏 */
  --sidebar: var(--background);
  --sidebar-foreground: var(--foreground);
  --sidebar-primary: var(--primary);
  --sidebar-primary-foreground: var(--primary-foreground);
  --sidebar-accent: var(--accent);
  --sidebar-accent-foreground: var(--accent-foreground);
  --sidebar-border: var(--border);
  --sidebar-ring: var(--ring);
}
```

### 暗色模式色彩变量

```css
.dark {
  /* 背景和前景反转 */
  --background: oklch(0.129 0.042 264.695);      /* 深紫灰 */
  --foreground: oklch(0.984 0.003 247.858);      /* 浅灰白 */

  /* 卡片更深 */
  --card: oklch(0.14 0.04 259.21);              /* 深色卡片 */
  --card-foreground: oklch(0.984 0.003 247.858);

  /* 弹出层 */
  --popover: oklch(0.208 0.042 265.755);        /* 中紫灰 */
  --popover-foreground: oklch(0.984 0.003 247.858);

  /* 主色调反转 */
  --primary: oklch(0.929 0.013 255.508);        /* 浅色 */
  --primary-foreground: oklch(0.208 0.042 265.755);

  /* 次要色和柔和色 */
  --secondary: oklch(0.279 0.041 260.031);      /* 深灰紫 */
  --secondary-foreground: oklch(0.984 0.003 247.858);
  --muted: oklch(0.279 0.041 260.031);
  --muted-foreground: oklch(0.704 0.04 256.788);

  /* 强调色 */
  --accent: oklch(0.279 0.041 260.031);
  --accent-foreground: oklch(0.984 0.003 247.858);

  /* 危险色更亮 */
  --destructive: oklch(0.704 0.191 22.216);     /* 亮红橙 */

  /* 边框半透明 */
  --border: oklch(1 0 0 / 10%);                 /* 10% 白色 */
  --input: oklch(1 0 0 / 15%);                  /* 15% 白色 */
  --ring: oklch(0.551 0.027 264.364);           /* 深紫环 */

  /* 暗色模式图表 */
  --chart-1: oklch(0.488 0.243 264.376);        /* 紫色 */
  --chart-2: oklch(0.696 0.17 162.48);          /* 青绿色 */
  --chart-3: oklch(0.769 0.188 70.08);          /* 绿色 */
  --chart-4: oklch(0.627 0.265 303.9);          /* 紫粉色 */
  --chart-5: oklch(0.645 0.246 16.439);         /* 橙色 */
}
```

### Tailwind 配置

**文件**: `frontend/tailwind.config.js`

```js
/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ['class'], // 支持暗色模式类切换
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
        secondary: {
          DEFAULT: 'hsl(var(--secondary))',
          foreground: 'hsl(var(--secondary-foreground))',
        },
        muted: {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        accent: {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
        },
        destructive: {
          DEFAULT: 'hsl(var(--destructive))',
          foreground: 'hsl(var(--destructive-foreground))',
        },
        card: {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
        },
        popover: {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
      },
      borderRadius: {
        lg: 'var(--radius)',
        md: 'calc(var(--radius) - 2px)',
        sm: 'calc(var(--radius) - 4px)',
      },
    },
  },
  plugins: [],
}
```

---

## A.2 组件别名配置

**文件**: `frontend/components.json` (shadcn/ui 配置)

```json
{
  "$schema": "https://ui.shadcn.com/schema.json",
  "style": "new-york",
  "rsc": false,
  "tsx": true,
  "tailwind": {
    "config": "tailwind.config.js",
    "css": "src/styles/index.css",
    "baseColor": "slate",
    "cssVariables": true
  },
  "aliases": {
    "components": "@/components",
    "utils": "@/lib/utils",
    "ui": "@/components/ui",
    "lib": "@/lib",
    "hooks": "@/hooks"
  },
  "iconLibrary": "lucide"
}
```

**路径映射** - `tsconfig.json`:

```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"],
      "@/components/*": ["./src/components/*"],
      "@/lib/*": ["./src/lib/*"],
      "@/hooks/*": ["./src/hooks/*"],
      "@/styles/*": ["./src/styles/*"],
      "@/assets/*": ["./src/assets/*"]
    }
  }
}
```

---

## A.3 基础样式文件

**文件**: `frontend/src/styles/index.css`

```css
@import 'tailwindcss';

@layer base {
  * {
    @apply border-border outline-ring/50;
    scrollbar-width: thin;
    scrollbar-color: var(--border) transparent;
  }

  html {
    @apply overflow-x-hidden;
  }

  body {
    @apply min-h-svh w-full bg-background text-foreground;
  }

  /* 按钮鼠标样式 */
  button:not(:disabled),
  [role='button']:not(:disabled) {
    cursor: pointer;
  }

  /* 防止移动端输入框缩放 */
  @media screen and (max-width: 767px) {
    input,
    select,
    textarea {
      font-size: 16px !important;
    }
  }
}
```

---

# 附录 B: UI 组件详细设计

## B.1 认证页面布局

### AuthLayout 组件

**文件**: `frontend/src/components/layout/AuthLayout.tsx`

```tsx
import { ReactNode } from 'react'

interface AuthLayoutProps {
  children: ReactNode
}

export function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className="container grid h-svh max-w-none items-center justify-center">
      <div className="mx-auto flex w-full flex-col justify-center space-y-2 py-8 sm:w-[480px] sm:p-8">
        {children}
      </div>
    </div>
  )
}
```

**布局说明**:
- `h-svh`: 全屏高度 (100svh - 视口高度，排除移动端浏览器 UI)
- `container`: 内容居中容器
- `grid`: 网格布局实现垂直居中
- `sm:w-[480px]`: 小屏幕以上限制宽度为 480px
- `sm:p-8`: 小屏幕以上添加内边距

### 登录/注册卡片结构

```
┌─────────────────────────────────────┐
│         [Logo] Artisan Coder         │ ← Logo 区域
├─────────────────────────────────────┤
│ Card                                │
│ ├─ CardHeader                       │
│ │  ├─ CardTitle: "Sign in"          │
│ │  └─ CardDescription               │
│ ├─ CardContent                      │
│ │  └─ Form                          │
│ │     ├─ Email Input                │
│ │     ├─ Password Input             │
│ │     ├─ Sign In Button             │
│ │     ├─ Divider: "Or continue with"│
│ │     └─ Social Buttons (GitHub)    │
│ └─ CardFooter                       │
│    └─ Terms & Privacy Links         │
└─────────────────────────────────────┘
```

---

## B.2 表单组件设计

### 登录表单

**文件**: `frontend/src/components/auth/LoginForm.tsx`

**字段**:
1. **邮箱** (`email`)
   - 类型: `email`
   - 验证: 必填，邮箱格式
   - 占位符: `name@example.com`
   - 错误提示: "Please enter your email"

2. **密码** (`password`)
   - 类型: `password` (带眼睛图标切换)
   - 验证: 必填，最少 7 个字符
   - 占位符: `********`
   - 右上角链接: "Forgot password?"

3. **按钮**
   - 主按钮: "Sign in" (带 `LogIn` 图标)
   - Loading 状态: 替换为 `Loader2` 旋转图标
   - 社交登录: GitHub 按钮

**表单验证**:
```tsx
import { z } from 'zod'

export const loginSchema = z.object({
  email: z
    .string()
    .min(1, 'Please enter your email')
    .email('Invalid email address'),
  password: z
    .string()
    .min(1, 'Please enter your password')
    .min(7, 'Password must be at least 7 characters long'),
})

export type LoginFormData = z.infer<typeof loginSchema>
```

**表单样式**:
```tsx
<form className="grid gap-3">
  {/* 表单字段间距 gap-3 */}

  {/* 分隔线 */}
  <div className="relative my-2">
    <div className="absolute inset-0 flex items-center">
      <span className="w-full border-t" />
    </div>
    <div className="relative flex justify-center text-xs uppercase">
      <span className="bg-background px-2 text-muted-foreground">
        Or continue with
      </span>
    </div>
  </div>

  {/* 社交按钮 */}
  <div className="grid grid-cols-2 gap-2">
    <Button variant="outline" type="button">
      <IconGithub className="h-4 w-4" /> GitHub
    </Button>
  </div>
</form>
```

### 注册表单

**文件**: `frontend/src/components/auth/RegisterForm.tsx`

**额外字段**:
1. **用户名** (`username`)
   - 验证: 必填，3-20 字符
   - 占位符: `johndoe`

2. **确认密码** (`confirmPassword`)
   - 验证: 必填，与密码一致
   - 占位符: `********`

**表单验证**:
```tsx
export const registerSchema = z
  .object({
    username: z
      .string()
      .min(1, 'Please enter your username')
      .min(3, 'Username must be at least 3 characters')
      .max(20, 'Username must not exceed 20 characters')
      .regex(/^[a-zA-Z0-9_-]+$/, 'Username can only contain letters, numbers, hyphens, and underscores'),
    email: z
      .string()
      .min(1, 'Please enter your email')
      .email('Invalid email address'),
    password: z
      .string()
      .min(1, 'Please enter your password')
      .min(7, 'Password must be at least 7 characters long'),
    confirmPassword: z
      .string()
      .min(1, 'Please confirm your password'),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match.",
    path: ['confirmPassword'],
  })

export type RegisterFormData = z.infer<typeof registerSchema>
```

### 密码输入组件

**文件**: `frontend/src/components/ui/password-input.tsx`

```tsx
import { useState } from 'react'
import { Eye, EyeOff } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

interface PasswordInputProps extends React.ComponentProps<typeof Input> {}

export function PasswordInput({ className, ...props }: PasswordInputProps) {
  const [showPassword, setShowPassword] = useState(false)

  return (
    <div className="relative">
      <Input
        type={showPassword ? 'text' : 'password'}
        className={cn('pr-10', className)}
        {...props}
      />
      <button
        type="button"
        onClick={() => setShowPassword(!showPassword)}
        className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
      >
        {showPassword ? (
          <EyeOff className="h-4 w-4" />
        ) : (
          <Eye className="h-4 w-4" />
        )}
      </button>
    </div>
  )
}
```

---

## B.3 Dashboard 布局设计

### AuthenticatedLayout 结构

```
┌─────────────────────────────────────────────────────────┐
│ SearchProvider                                          │
│ ┌─ LayoutProvider ──────────────────────────────────┐   │
│ │ ┌─ SidebarProvider ─────────────────────────────┐ │   │
│ │ │ ┌─────────┐ ┌─────────────────────────────┐   │ │   │
│ │ │ │         │ │  Header (固定顶部)           │   │ │   │
│ │ │ │Sidebar  │ ├─────────────────────────────┤   │ │   │
│ │ │ │         │ │                             │   │ │   │
│ │ │ │(可折叠) │ │   Main 内容区域              │   │ │   │
│ │ │ │         │ │   (可滚动)                   │   │ │   │
│ │ │ │         │ │                             │   │ │   │
│ │ │ │         │ │                             │   │ │   │
│ │ │ └─────────┘ └─────────────────────────────┘   │ │   │
│ │ └─────────────────────────────────────────────┘ │   │
│ └───────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

### Header 组件

**文件**: `frontend/src/components/layout/Header.tsx`

```tsx
import { ReactNode, useEffect, useState } from 'react'
import { cn } from '@/lib/utils'
import { Separator } from '@/components/ui/separator'
import { SidebarTrigger } from '@/components/ui/sidebar'
import { ThemeSwitch } from '@/components/theme-switch'
import { UserMenu } from '@/components/user-menu'

interface HeaderProps {
  children?: ReactNode
}

export function Header({ children }: HeaderProps) {
  const [offset, setOffset] = useState(0)

  useEffect(() => {
    const onScroll = () => {
      setOffset(document.body.scrollTop || document.documentElement.scrollTop)
    }
    document.addEventListener('scroll', onScroll, { passive: true })
    return () => document.removeEventListener('scroll', onScroll)
  }, [])

  return (
    <header
      className={cn(
        'z-50 h-16 sticky top-0',
        offset > 10 && 'shadow'
      )}
    >
      <div className="flex h-full items-center gap-3 p-4">
        <SidebarTrigger variant="outline" />
        <Separator orientation="vertical" className="h-6" />
        {children}
        <div className="ms-auto flex items-center gap-2">
          <ThemeSwitch />
          <UserMenu />
        </div>
      </div>
    </header>
  )
}
```

### Sidebar 组件

**文件**: `frontend/src/components/layout/Sidebar.tsx`

```tsx
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '@/components/ui/sidebar'
import { AppLogo } from '@/components/app-logo'
import { NavGroup } from '@/components/nav-group'
import { NavUser } from '@/components/nav-user'

export function AppSidebar() {
  return (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <AppLogo />
      </SidebarHeader>
      <SidebarContent>
        <NavGroup
          title="General"
          items={[
            { title: 'Dashboard', url: '/', icon: LayoutDashboard },
            { title: 'Projects', url: '/projects', icon: FolderKanban },
            { title: 'Settings', url: '/settings', icon: Settings },
          ]}
        />
      </SidebarContent>
      <SidebarFooter>
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
```

### Main 组件

**文件**: `frontend/src/components/layout/Main.tsx`

```tsx
import { ReactNode } from 'react'
import { cn } from '@/lib/utils'

interface MainProps {
  children: ReactNode
  className?: string
}

export function Main({ children, className }: MainProps) {
  return (
    <main className={cn(
      'px-4 py-6',
      'max-w-7xl mx-auto w-full',
      className
    )}>
      {children}
    </main>
  )
}
```

---

## B.4 Dashboard 内容区域

### 页面标题区

```tsx
<div className="mb-2 flex items-center justify-between">
  <h1 className="text-2xl font-bold tracking-tight">Dashboard</h1>
  <div className="flex items-center gap-2">
    <Button>Create Project</Button>
  </div>
</div>
```

### 统计卡片网格

```tsx
<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
  <StatsCard
    title="Total Projects"
    value={stats.total}
    icon={<FolderKanban className="h-4 w-4" />}
    description="+2 from last month"
  />
  <StatsCard
    title="Running"
    value={stats.running}
    icon={<Play className="h-4 w-4" />}
    description="Active containers"
  />
  <StatsCard
    title="Stopped"
    value={stats.stopped}
    icon={<StopCircle className="h-4 w-4" />}
  />
  <StatsCard
    title="Errors"
    value={stats.error}
    icon={<AlertCircle className="h-4 w-4 text-destructive" />}
  />
</div>
```

### StatsCard 组件

**文件**: `frontend/src/components/dashboard/StatsCard.tsx`

```tsx
import { ReactNode } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { cn } from '@/lib/utils'

interface StatsCardProps {
  title: string
  value: number
  icon: ReactNode
  description?: string
  className?: string
}

export function StatsCard({ title, value, icon, description, className }: StatsCardProps) {
  return (
    <Card className={cn('', className)}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">{value}</div>
        {description && (
          <p className="text-xs text-muted-foreground">{description}</p>
        )}
      </CardContent>
    </Card>
  )
}
```

### 项目卡片网格

```tsx
<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
  {projects.map((project) => (
    <ProjectCard key={project.id} project={project} />
  ))}
</div>
```

### ProjectCard 组件

**文件**: `frontend/src/components/project/ProjectCard.tsx`

```tsx
import { Project } from '@/types/models'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { MoreVertical, Play, Stop, ExternalLink, Trash2 } from 'lucide-react'

interface ProjectCardProps {
  project: Project
  onStart?: () => void
  onStop?: () => void
  onOpen?: () => void
  onDelete?: () => void
}

export function ProjectCard({ project, onStart, onStop, onOpen, onDelete }: ProjectCardProps) {
  const statusColors = {
    running: 'bg-green-500/10 text-green-700 dark:text-green-400',
    stopped: 'bg-gray-500/10 text-gray-700 dark:text-gray-400',
    error: 'bg-red-500/10 text-red-700 dark:text-red-400',
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <CardTitle className="text-lg">{project.name}</CardTitle>
          <Badge variant="outline" className={statusColors[project.status]}>
            {project.status}
          </Badge>
        </div>
        {project.description && (
          <p className="text-sm text-muted-foreground">{project.description}</p>
        )}
      </CardHeader>
      <CardContent>
        <div className="space-y-2 text-sm">
          <div className="flex items-center gap-2">
            <span className="text-muted-foreground">Type:</span>
            <span className="capitalize">{project.projectType}</span>
          </div>
          {project.gitRepo && (
            <div className="flex items-center gap-2 truncate">
              <span className="text-muted-foreground">Repo:</span>
              <span className="truncate">{project.gitRepo}</span>
            </div>
          )}
        </div>
      </CardContent>
      <CardFooter className="flex justify-between">
        <Button
          variant="default"
          size="sm"
          onClick={project.status === 'running' ? onStop : onStart}
        >
          {project.status === 'running' ? (
            <><Stop className="h-4 w-4 mr-1" /> Stop</>
          ) : (
            <><Play className="h-4 w-4 mr-1" /> Start</>
          )}
        </Button>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon">
              <MoreVertical className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={onOpen}>
              <ExternalLink className="h-4 w-4 mr-2" /> Open
            </DropdownMenuItem>
            <DropdownMenuItem onClick={onDelete} className="text-destructive">
              <Trash2 className="h-4 w-4 mr-2" /> Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </CardFooter>
    </Card>
  )
}
```

### EmptyState 组件

**文件**: `frontend/src/components/dashboard/EmptyState.tsx`

```tsx
import { ReactNode } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'

interface EmptyStateProps {
  icon?: ReactNode
  title: string
  description?: string
  actionLabel?: string
  onAction?: () => void
}

export function EmptyState({ icon, title, description, actionLabel, onAction }: EmptyStateProps) {
  return (
    <Card>
      <CardContent className="flex flex-col items-center justify-center py-12">
        {icon && <div className="mb-4 text-muted-foreground">{icon}</div>}
        <h3 className="text-lg font-semibold">{title}</h3>
        {description && (
          <p className="text-sm text-muted-foreground mb-4">{description}</p>
        )}
        {actionLabel && onAction && (
          <Button onClick={onAction}>{actionLabel}</Button>
        )}
      </CardContent>
    </Card>
  )
}
```

---

# 附录 C: 技术栈对比

## C.1 shadcn-admin vs Artisan Coder

| 类别 | shadcn-admin | Artisan Coder |
|------|--------------|---------------|
| **路由** | @tanstack/react-router | react-router-dom |
| **表单** | react-hook-form + zod | react-hook-form + zod |
| **状态管理** | zustand | zustand |
| **HTTP 客户端** | axios | axios |
| **UI 组件** | shadcn/ui | shadcn/ui |
| **图标** | lucide-react | lucide-react |
| **日期处理** | date-fns | date-fns |
| **Toast** | sonner | sonner |
| **Markdown** | - | react-markdown |
| **代码高亮** | - | react-syntax-highlighter |

## C.2 依赖版本参考

```json
{
  "dependencies": {
    "react": "^19.x",
    "react-dom": "^19.x",
    "react-router-dom": "^7.x",
    "zustand": "^5.x",
    "axios": "^1.x",
    "react-hook-form": "^7.x",
    "zod": "^4.x",
    "@hookform/resolvers": "^5.x",
    "tailwindcss": "^4.x",
    "class-variance-authority": "^0.7.x",
    "clsx": "^2.x",
    "tailwind-merge": "^3.x",
    "lucide-react": "^0.x",
    "date-fns": "^4.x",
    "sonner": "^2.x",
    "react-markdown": "^9.x",
    "react-syntax-highlighter": "^15.x"
  },
  "devDependencies": {
    "@types/react": "^19.x",
    "@types/react-dom": "^19.x",
    "@types/react-syntax-highlighter": "^15.x",
    "@vitejs/plugin-react": "^4.x",
    "typescript": "~5.9.x",
    "vite": "^6.x",
    "prettier": "^3.x",
    "eslint": "^9.x"
  }
}
```

---

# 附录 D: 路由配置示例

**文件**: `frontend/src/router.tsx`

```tsx
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { Login, Register } from '@/pages/auth'
import { Dashboard } from '@/pages/dashboard'
import { ProtectedRoute, PublicRoute } from '@/components/auth'
import { AuthenticatedLayout } from '@/components/layout'

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
        <AuthenticatedLayout />
      </ProtectedRoute>
    ),
    children: [
      {
        index: true,
        element: <Dashboard />,
      },
      {
        path: 'projects',
        element: <div>Projects Page</div>, // Phase 2
      },
      {
        path: 'settings',
        element: <div>Settings Page</div>, // Phase 2
      },
    ],
  },
])

export function AppRouter() {
  return <RouterProvider router={router} />
}
```

---

# 附录 E: 响应式断点

```css
/* Tailwind 默认断点 */
sm: 640px   /* 小型设备 */
md: 768px   /* 平板设备 */
lg: 1024px  /* 桌面设备 */
xl: 1280px  /* 大型桌面 */
2xl: 1536px /* 超大屏幕 */
```

**响应式布局示例**:
```tsx
{/* 统计卡片 */}
<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
  {/* 移动端: 1列, 小屏: 2列, 大屏: 4列 */}
</div>

{/* 项目卡片 */}
<div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
  {/* 移动端: 1列, 小屏: 2列, 大屏: 3列, 超大: 4列 */}
</div>
```

---

# 附录 F: Toast 通知配置

**使用 sonner**:

```tsx
// src/main.tsx
import { Toaster } from 'sonner'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
    <Toaster
      position="top-right"
      expand={false}
      richColors
      closeButton
    />
  </React.StrictMode>
)
```

**使用示例**:
```tsx
import { toast } from 'sonner'

// 成功提示
toast.success('登录成功')

// 错误提示
toast.error('登录失败', {
  description: '用户名或密码错误'
})

// 加载提示
toast.promise(loginPromise, {
  loading: '登录中...',
  success: '登录成功',
  error: '登录失败'
})
```

---

# 附录 G: 图标库 (lucide-react)

**常用图标**:
```tsx
import {
  // 布局
  LayoutDashboard,
  Sidebar,

  // 项目
  FolderKanban,
  FolderOpen,
  Plus,

  // 状态
  Play,
  StopCircle,
  Pause,
  AlertCircle,
  CheckCircle2,

  // 导航
  Settings,
  User,
  LogOut,
  Menu,

  // 认证
  LogIn,
  UserPlus,
  Eye,
  EyeOff,

  // 其他
  MoreVertical,
  ExternalLink,
  Trash2,
  Search,
  Bell,
} from 'lucide-react'
```

---

**文档版本**: v1.0
**最后更新**: 2025-12-26
**参考来源**: shadcn-admin 项目
