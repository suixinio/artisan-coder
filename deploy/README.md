# Artisan Coder 数据库部署

本目录包含 Artisan Coder 项目的 PostgreSQL 数据库 Docker Compose 配置。

## 快速开始

### 1. 配置环境变量

```bash
# 复制示例配置文件
cp .env.example .env

# 编辑配置文件，修改数据库密码等敏感信息
vim .env
```

### 2. 启动数据库

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f postgres
```

### 3. 验证服务

```bash
# 检查容器状态
docker-compose ps

# 测试数据库连接
docker exec -it artisan-coder-postgres psql -U artisan -d artisan_coder
```

### 4. 连接数据库

- **主机**: localhost
- **端口**: 5433
- **数据库**: artisan_coder
- **用户名**: artisan
- **密码**: (在 .env 文件中配置)

## 常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 停止服务并删除数据卷（⚠️ 会删除所有数据）
docker-compose down -v

# 重启服务
docker-compose restart

# 查看日志
docker-compose logs -f postgres

# 进入数据库
docker exec -it artisan-coder-postgres psql -U artisan -d artisan_coder

# 备份数据库
docker exec artisan-coder-postgres pg_dump -U artisan artisan_coder > backup.sql

# 恢复数据库
docker exec -i artisan-coder-postgres psql -U artisan artisan_coder < backup.sql
```

## 配置说明

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `POSTGRES_USER` | 数据库用户名 | `artisan` |
| `POSTGRES_PASSWORD` | 数据库密码（必填） | - |
| `POSTGRES_DB` | 数据库名称 | `artisan_coder` |
| `TZ` | 时区 | `Asia/Shanghai` |

### 端口映射

- 宿主机端口 `5433` 映射到容器端口 `5432`
- 使用 5433 是为了避免与宿主机可能已有的 PostgreSQL 服务冲突

### 数据持久化

- 数据存储在 Docker 命名卷 `postgres_data` 中
- 即使容器被删除，数据也会保留
- 使用 `docker-compose down -v` 可以删除数据卷

## 生产环境建议

1. **修改默认密码**: 务必修改 `.env` 文件中的默认密码
2. **备份策略**: 定期备份数据库
3. **网络隔离**: 在生产环境中，不要将数据库端口暴露到公网
4. **资源限制**: 根据需要添加 CPU 和内存限制
5. **日志管理**: 配置日志轮转策略

## 故障排查

### 容器无法启动

```bash
# 查看详细日志
docker-compose logs postgres

# 检查端口占用
netstat -tuln | grep 5433
```

### 无法连接数据库

```bash
# 检查容器状态
docker-compose ps

# 检查健康状态
docker inspect artisan-coder-postgres | grep -A 10 Health

# 测试连接
docker exec -it artisan-coder-postgres pg_isready -U artisan
```

### 数据损坏

```bash
# 如果数据损坏，可以重建容器（会丢失数据）
docker-compose down -v
docker-compose up -d
```

## License

MIT
