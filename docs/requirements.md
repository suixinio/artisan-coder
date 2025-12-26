# Artisan Coder 需求文档

## 1. 项目概述

Artisan Coder 是一个通用的 AI 编码助手平台，采用三层架构设计，为开发者提供智能的代码辅助和项目管理能力。

### 1.1 核心定位

- **多 Agent 支持**: 通过 ACP (Agent Communication Protocol) 协议支持多种 AI 编码助手
- **容器化隔离**: 每个项目对应独立的容器环境，确保安全和隔离
- **全栈解决方案**: 提供 Web 端、API 端、CLI 端完整的解决方案

## 2. 架构设计

### 2.1 三端架构

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    Web 端       │────▶│    API 端       │────▶│    CLI 端       │
│  (用户界面)      │     │  (业务逻辑)      │     │  (守护进程)      │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                                                        │
                                                        ▼
                                              ┌─────────────────┐
                                              │  AI Agent       │
                                              │ (OpenCode/      │
                                              │  Claude Code/   │
                                              │  Codex/...)     │
                                              └─────────────────┘
```

### 2.2 各端职责

#### Web 端 (前端)
- **技术栈**: React 18+ + TypeScript + Vite + shadcn/ui
- **职责**:
  - 用户登录注册
  - 项目管理界面
  - 会话聊天界面
  - 文件浏览器
  - SSH 终端
  - 实时通信 (WebSocket)

#### API 端 (后端)
- **技术栈**: Go + Gin + GORM + Docker SDK + PostgreSQL
- **职责**:
  - 用户认证授权
  - 项目管理 (CRUD)
  - 容器生命周期管理
  - 会话管理
  - 消息路由
  - WebSocket 连接管理
  - 与 CLI 端通信

#### CLI 端 (容器守护进程)
- **技术栈**: Go + WebSocket + Docker
- **职责**:
  - 运行在容器中的守护进程
  - 实现 ACP 协议
  - 与 AI Agent 交互
  - 文件系统操作
  - 命令执行
  - 与 API 端双向通信

## 3. 核心功能需求

### 3.1 用户认证

| 功能 | 描述 |
|------|------|
| 用户注册 | 邮箱/用户名注册，密码加密存储 |
| 用户登录 | JWT Token 认证，支持刷新 Token |
| 登出 | 清除客户端 Token，服务端会话失效 |
| 权限管理 | 基于用户的数据隔离 |

### 3.2 项目管理 (Project)

| 功能 | 描述 |
|------|------|
| 创建项目 | 支持空项目或关联 GitHub 仓库 |
| 项目列表 | 展示用户所有项目，支持搜索/过滤 |
| 项目详情 | 查看项目信息、容器状态、文件目录 |
| 启动/停止 | 控制项目容器运行状态 |
| 删除项目 | 删除项目及相关容器和数据 |
| SSH 连接 | Web terminal 直接连入容器 |

#### 项目类型

1. **空项目**: 从零开始的开发环境
2. **仓库项目**: 关联 GitHub 仓库，自动克隆代码

#### 项目状态

- `stopped`: 容器已停止
- `running`: 容器运行中
- `error`: 容器异常

### 3.3 会话管理 (Session)

| 功能 | 描述 |
|------|------|
| 创建会话 | 在项目下创建新的 AI 对话会话 |
| 会话列表 | 展示项目下所有会话 |
| 选择 Agent | 为会话选择 AI Agent 类型 |
| 切换 Agent | 同一会话可切换不同 Agent |
| 删除会话 | 删除会话及相关消息历史 |

#### 支持的 Agent 类型

- **OpenCode**: 开源代码助手
- **Claude Code**: Anthropic Claude 编码助手
- **Codex**: OpenAI Codex
- **其他**: 通过 ACP 协议可扩展更多 Agent

### 3.4 消息聊天 (Chat)

| 功能 | 描述 |
|------|------|
| 发送消息 | 用户向 AI 发送消息 |
| 流式响应 | AI 回复实时流式展示 |
| 多轮对话 | 保存完整对话历史 |
| 工具调用展示 | 显示 AI 使用的工具和操作 |
| 代码高亮 | Markdown 代码块语法高亮 |
| 消息操作 | 复制、删除、重新生成 |

### 3.5 文件浏览 (File Explorer)

| 功能 | 描述 |
|------|------|
| 目录树 | 展示容器内文件目录结构 |
| 文件预览 | 查看文件内容（代码/文本） |
| 文件编辑 | 在线编辑文件（可选） |
| 文件上传 | 上传文件到容器 |
| 文件下载 | 从容器下载文件 |

### 3.6 SSH 终端

| 功能 | 描述 |
|------|------|
| Web Terminal | 浏览器内嵌终端 |
| Shell 访问 | 直接连入容器 Shell |
| 多标签页 | 支持多个终端会话 |

## 4. 数据模型

### 4.1 User (用户)

```typescript
interface User {
  id: string;          // UUID
  username: string;    // 用户名
  email: string;       // 邮箱
  passwordHash: string; // 密码哈希
  createdAt: Date;
  updatedAt: Date;
}
```

### 4.2 Project (项目)

```typescript
interface Project {
  id: string;          // UUID
  userId: string;      // 所属用户
  name: string;        // 项目名称
  description?: string;// 项目描述
  projectType: 'empty' | 'repo'; // 项目类型
  gitRepo?: string;    // Git 仓库地址
  gitBranch?: string;  // Git 分支
  containerId?: string;// Docker 容器 ID
  containerName?: string; // 容器名称
  status: 'stopped' | 'running' | 'error'; // 状态
  sshPort?: number;    // SSH 映射端口
  createdAt: Date;
  updatedAt: Date;
}
```

### 4.3 Session (会话)

```typescript
interface Session {
  id: string;          // UUID
  projectId: string;   // 所属项目
  name: string;        // 会话名称
  agentType: string;   // Agent 类型 (opencode/claude-code/codex/...)
  agentConfig?: Record<string, any>; // Agent 配置
  createdAt: Date;
  updatedAt: Date;
}
```

### 4.4 Message (消息)

```typescript
interface Message {
  id: string;          // UUID
  sessionId: string;   // 所属会话
  role: 'user' | 'assistant' | 'system';
  content: string;     // 消息内容 (Markdown)
  metadata?: {
    model?: string;    // 使用的模型
    tokensUsed?: number; // Token 消耗
    toolCalls?: ToolCall[]; // 工具调用
    thinking?: string; // 思考过程 (可选)
  };
  createdAt: Date;
}
```

### 4.5 ToolCall (工具调用)

```typescript
interface ToolCall {
  id: string;
  name: string;        // 工具名称
  args: any;           // 工具参数
  result?: any;        // 工具结果
  status: 'pending' | 'success' | 'error';
}
```

## 5. API 接口

### 5.1 认证接口

```
POST   /api/auth/register  # 用户注册
POST   /api/auth/login     # 用户登录
POST   /api/auth/logout    # 用户登出
POST   /api/auth/refresh   # 刷新 Token
GET    /api/auth/me        # 获取当前用户信息
```

### 5.2 项目接口

```
GET    /api/projects           # 获取项目列表
POST   /api/projects           # 创建项目
GET    /api/projects/:id       # 获取项目详情
PUT    /api/projects/:id       # 更新项目
DELETE /api/projects/:id       # 删除项目
POST   /api/projects/:id/start # 启动项目容器
POST   /api/projects/:id/stop  # 停止项目容器
GET    /api/projects/:id/files # 获取文件列表
```

### 5.3 会话接口

```
GET    /api/projects/:projectId/sessions      # 获取会话列表
POST   /api/projects/:projectId/sessions      # 创建会话
GET    /api/sessions/:id                      # 获取会话详情
DELETE /api/sessions/:id                      # 删除会话
PUT    /api/sessions/:id/agent                # 切换 Agent
```

### 5.4 消息接口

```
GET    /api/sessions/:sessionId/messages  # 获取消息历史
POST   /api/sessions/:sessionId/messages  # 发送消息
DELETE /api/messages/:id                 # 删除消息
```

### 5.5 WebSocket

```
WS     /api/ws/chat?sessionId=:id     # 聊天流式通信
WS     /api/ws/terminal?projectId=:id  # SSH 终端
WS     /api/ws/files?projectId=:id    # 文件变更通知
```

## 6. 通信流程

### 6.1 创建项目流程

```
1. 用户在 Web 端创建项目
2. API 端创建项目记录（数据库）
3. API 端调用 Docker API 创建容器
4. CLI 端在容器中启动守护进程
5. CLI 端通过 WebSocket 连接 API 端
6. 容器状态更新为 running
```

### 6.2 聊天消息流程

```
1. 用户在 Web 端发送消息
2. Web 端通过 WebSocket 发送到 API 端
3. API 端保存消息到数据库
4. API 端转发消息到 CLI 端（WebSocket）
5. CLI 端通过 ACP 协议发送给 AI Agent
6. AI Agent 流式响应
7. CLI 端转发响应到 API 端
8. API 端流式推送到 Web 端
9. Web 端实时展示 AI 回复
```

### 6.3 ACP 协议 (Agent Communication Protocol)

CLI 端与 AI Agent 之间的通信协议：

```typescript
// ACP 消息格式
interface ACPMessage {
  version: string;      // 协议版本
  type: 'chat' | 'tool' | 'file' | 'system';
  payload: any;
}

// Chat 类型
interface ACPChatMessage {
  role: 'user' | 'assistant';
  content: string;
  stream?: boolean;     // 是否流式
}

// Tool 类型
interface ACPToolMessage {
  tool: string;
  action: 'call' | 'result';
  params?: any;
  result?: any;
}
```

## 7. 容器管理

### 7.1 容器配置

```yaml
BaseImage: ubuntu:22.04
Runtime: nodejs / python / golang (可选)
InstalledPackages:
  - git
  - curl
  - vim
  - openssh-server
Mounts:
  - /workspace
Ports:
  - SSH: 映射到主机随机端口
  - CLI: 内部通信端口
Environment:
  - CLI_ENDPOINT: ws://api:8080/ws/cli
  - AGENT_TYPE: configurable
```

### 7.2 CLI 守护进程

```go
// CLI 端核心职责
type CLIDaemon struct {
    AgentAdapter AgentAdapter  // Agent 适配器
    ACPClient    *ACPClient    // ACP 协议客户端
    APIClient    *APIClient    // API 端通信
    FileMgr      *FileManager  // 文件管理
}

// 启动流程
1. 连接 API 端 WebSocket
2. 初始化 ACP 协议
3. 加载 Agent 配置
4. 等待消息处理
```

## 8. 安全考虑

### 8.1 认证授权
- JWT Token 认证
- Token 过期自动刷新
- 资源级权限控制 (用户只能访问自己的项目)

### 8.2 容器隔离
- 每个项目独立容器
- 资源限制 (CPU/内存)
- 网络隔离

### 8.3 数据安全
- 密码哈希存储 (bcrypt)
- HTTPS/WSS 加密通信
- SQL 注入防护

### 8.4 输入验证
- 用户输入验证和清理
- 文件操作权限检查
- 命令注入防护

## 9. 非功能性需求

### 9.1 性能
- 消息响应时间 < 500ms (不含 AI 处理)
- 支持并发会话数 > 100
- WebSocket 长连接稳定性

### 9.2 可扩展性
- 水平扩展 API 端
- 支持分布式部署
- Agent 插件化

### 9.3 可维护性
- 模块化架构
- 统一日志规范
- 健康检查接口

### 9.4 可用性
- 服务可用性 > 99%
- 容器故障自动恢复
- 数据库持久化

## 10. 技术选型

### 10.1 Web 端
- React 18+: UI 框架
- TypeScript: 类型安全
- Vite: 构建工具
- shadcn/ui: UI 组件库
- Zustand: 状态管理
- React Router: 路由
- Axios: HTTP 客户端
- react-markdown: Markdown 渲染

### 10.2 API 端
- Go: 高性能后端
- Gin: Web 框架
- GORM: ORM
- PostgreSQL: 关系型数据库
- Docker SDK: 容器管理
- JWT: 认证
- WebSocket: 实时通信

### 10.3 CLI 端
- Go: 守护进程
- WebSocket: 与 API 通信
- ACP: 自定义协议
- Docker: 容器运行时

## 11. 开发阶段

### Phase 1: 基础框架
- [ ] 项目脚手架搭建
- [ ] 数据库设计
- [ ] 基础认证系统
- [ ] 前端路由和布局

### Phase 2: 核心功能
- [ ] 项目管理功能
- [ ] 容器创建和管理
- [ ] CLI 守护进程
- [ ] ACP 协议实现

### Phase 3: AI 集成
- [ ] Agent 适配器
- [ ] 聊天功能
- [ ] 流式响应
- [ ] 工具调用展示

### Phase 4: 增强功能
- [ ] 文件浏览器
- [ ] SSH 终端
- [ ] 多 Agent 支持
- [ ] 会话管理优化

### Phase 5: 优化部署
- [ ] 性能优化
- [ ] 安全加固
- [ ] Docker Compose 部署
- [ ] 文档完善

## 12. 参考资料

- [shadcn/ui](https://ui.shadcn.com/)
- [React Router](https://reactrouter.com/)
- [Gin Framework](https://gin-gonic.com/)
- [Docker SDK for Go](https://docs.docker.com/engine/api/sdk/)
- [WebSocket](https://websocket.org/)
