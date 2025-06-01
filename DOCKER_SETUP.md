# 🐳 Docker 部署指南

## 📋 配置说明

`docker-compose.yml` 现在完全从 `.env` 文件读取配置，确保了安全性和灵活性。

## 🚀 快速开始

### 1. 准备环境文件

**开发环境:**
```bash
# .env 文件已经存在，包含开发环境配置
# 确认配置正确即可
```

**生产环境:**
```bash
# 复制生产环境模板
cp .env.production .env

# 编辑 .env 文件，填入真实的生产环境配置
nano .env
```

### 2. 启动服务

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 3. 访问服务

- **API服务**: http://localhost:3000
- **API文档**: http://localhost:3000/swagger/index.html
- **健康检查**: http://localhost:3000/api/health
- **数据库管理**: http://localhost:9080 (PgAdmin)

## 🔧 配置详情

### 环境变量读取
- **ar-backend服务**: 从 `.env` 文件读取所有配置
- **postgres服务**: 数据库配置从 `.env` 文件读取
- **特殊覆盖**: `BLUEPRINT_DB_HOST=postgres` (Docker内部网络)

### 端口映射
- **应用端口**: 3000 (HOST:CONTAINER)
- **数据库端口**: 5432 (HOST:CONTAINER)
- **PgAdmin端口**: 9080:80 (HOST:CONTAINER)

## 🗃️ 数据持久化

- **数据库数据**: `./pgdata/` 目录
- **PgAdmin配置**: `./pgadmin-data/` 目录

## 🔍 故障排除

### 查看服务状态
```bash
docker-compose ps
```

### 查看特定服务日志
```bash
docker-compose logs ar-backend
docker-compose logs postgres
```

### 重新构建服务
```bash
docker-compose up --build
```

### 清理并重启
```bash
docker-compose down
docker-compose up -d
```

## 🛡️ 安全注意事项

1. **生产环境**: 确保 `.env` 文件不被提交到版本控制
2. **密钥安全**: 使用强密码和真实的OAuth密钥
3. **网络安全**: 考虑使用反向代理(nginx)
4. **数据备份**: 定期备份 `pgdata` 目录

## 📞 支持

如有问题，请参考：
- [SECURITY.md](./SECURITY.md) - 安全配置指南
- [README.md](./README.md) - 项目总览 