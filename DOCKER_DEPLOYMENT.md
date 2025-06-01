# Docker 化部署指南

本项目提供高效的 Docker 化解决方案，后端容器化，前端通过构建脚本直接部署到 nginx。

## 🏗️ 架构概览

```
Internet → Nginx (SSL + 静态文件) → Backend container
                 ↓
              PostgreSQL
```

### 服务组件

- **nginx**: 反向代理、SSL 终端、前端静态文件服务
- **ar-backend**: Go 后端 API 服务 (容器化)
- **postgres**: PostgreSQL 数据库 (容器化)
- **certbot**: Let's Encrypt SSL 证书管理 (容器化)
- **前端**: React 应用构建后直接部署到 nginx 静态目录

## 🚀 快速开始

### 1. 环境准备

确保服务器上已安装：
- Docker (20.10+)
- Docker Compose (2.0+)
- Node.js (18+)
- npm
- curl
- 域名解析到服务器 IP

### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.docker .env.production

# 编辑配置文件
nano .env.production
```

重要配置项：
- `BLUEPRINT_DB_PASSWORD`: 数据库密码
- `JWT_SECRET`: JWT 密钥
- `GOOGLE_CLIENT_ID/SECRET`: Google OAuth 配置

### 3. 部署应用

```bash
# 添加执行权限
chmod +x scripts/*.sh

# 完整部署（推荐首次使用）
./scripts/deploy-docker.sh

# 仅更新前端
./scripts/deploy-docker.sh --frontend

# 仅更新后端
./scripts/deploy-docker.sh --backend

# 查看帮助
./scripts/deploy-docker.sh --help
```

### 4. 配置 SSL 证书

```bash
# 首次申请证书
./scripts/init-letsencrypt.sh
```

## 📂 目录结构

```
ar-backend/
├── docker-compose.yml          # 主要服务配置
├── Dockerfile                  # 后端镜像构建
├── client/
│   ├── dist/                   # 前端构建输出 (自动生成)
│   ├── src/                    # 前端源码
│   └── package.json            # 前端依赖配置
├── nginx/
│   └── nginx.conf              # nginx 配置
├── scripts/
│   ├── deploy-docker.sh        # 统一部署脚本 (合并版本)
│   └── init-letsencrypt.sh     # SSL 证书初始化
├── certbot/                    # SSL 证书存储
├── nginx/logs/                 # nginx 日志
└── backups/                    # 数据库备份
```

## 🔧 常用命令

### 部署和更新
```bash
# 完整部署（前端+后端+数据库）
./scripts/deploy-docker.sh

# 仅更新前端（快速）
./scripts/deploy-docker.sh --frontend

# 仅更新后端
./scripts/deploy-docker.sh --backend

# 查看使用帮助
./scripts/deploy-docker.sh --help
```

### Docker 服务管理
```bash
# 启动所有服务
docker-compose --env-file=.env.production up -d

# 启动特定服务
docker-compose --env-file=.env.production up -d nginx ar-backend

# 重启后端
docker-compose --env-file=.env.production restart ar-backend

# 重启 nginx
docker-compose --env-file=.env.production restart nginx

# 停止所有服务
docker-compose --env-file=.env.production down
```

### 查看状态和日志
```bash
# 查看所有容器状态
docker-compose --env-file=.env.production ps

# 查看日志
docker-compose --env-file=.env.production logs -f nginx
docker-compose --env-file=.env.production logs -f ar-backend
docker-compose --env-file=.env.production logs -f postgres
```

### 前端开发
```bash
# 进入前端目录
cd client

# 安装依赖
npm install

# 开发模式
npm run dev

# 构建生产版本
npm run build

# 返回根目录并部署
cd ..
./scripts/deploy-docker.sh --frontend
```

### 数据库管理
```bash
# 启动 pgAdmin（可选）
docker-compose --env-file=.env.production --profile admin up -d pgadmin

# 进入数据库容器
docker-compose --env-file=.env.production exec postgres psql -U ifoodme_user -d ifoodme_db

# 数据库备份
docker-compose --env-file=.env.production exec postgres pg_dump -U ifoodme_user ifoodme_db > backup.sql

# 数据库恢复
docker-compose --env-file=.env.production exec -T postgres psql -U ifoodme_user -d ifoodme_db < backup.sql
```

## 🔐 SSL 证书管理

### 自动续期
证书会每 12 小时自动检查并续期。

### 手动续期
```bash
docker-compose --env-file=.env.production run --rm certbot renew
docker-compose --env-file=.env.production exec nginx nginx -s reload
```

### 测试证书配置
```bash
# 检查证书状态
openssl s_client -connect www.ifoodme.com:443 -servername www.ifoodme.com

# 检查 SSL 评级
curl -s "https://api.ssllabs.com/api/v3/analyze?host=www.ifoodme.com"
```

## 🔍 监控和调试

### 健康检查
```bash
# 检查后端健康状态
curl -f https://www.ifoodme.com/api/health

# 检查前端
curl -f https://www.ifoodme.com/

# 检查所有容器健康状态
docker-compose --env-file=.env.production ps
```

### 前端问题排查
```bash
# 检查前端文件是否存在
ls -la client/dist/

# 检查 nginx 是否能访问前端文件
docker-compose --env-file=.env.production exec nginx ls -la /var/www/html/

# 重新构建前端
./scripts/deploy-docker.sh --frontend
```

### 后端问题排查
```bash
# 查看后端日志
docker-compose --env-file=.env.production logs --tail=100 ar-backend

# 进入后端容器
docker-compose --env-file=.env.production exec ar-backend /bin/sh

# 重新构建后端
./scripts/deploy-docker.sh --backend
```

## 🛠️ 故障排除

### 常见问题

1. **前端无法访问**
   ```bash
   # 检查前端是否构建
   ls client/dist/
   
   # 重新构建前端
   ./scripts/deploy-docker.sh --frontend
   
   # 检查 nginx 配置
   docker-compose --env-file=.env.production exec nginx nginx -t
   ```

2. **后端 API 无法访问**
   ```bash
   # 检查后端容器状态
   docker-compose --env-file=.env.production ps ar-backend
   
   # 查看后端日志
   docker-compose --env-file=.env.production logs ar-backend
   
   # 重新构建后端
   ./scripts/deploy-docker.sh --backend
   ```

3. **数据库连接失败**
   ```bash
   # 检查数据库是否启动
   docker-compose --env-file=.env.production exec postgres pg_isready
   
   # 检查网络连接
   docker-compose --env-file=.env.production exec ar-backend ping postgres
   ```

4. **SSL 证书问题**
   ```bash
   # 重新申请证书
   ./scripts/init-letsencrypt.sh
   
   # 检查证书文件
   ls -la certbot/letsencrypt/live/www.ifoodme.com/
   ```

## 🔄 更新部署

### 代码更新工作流
```bash
# 1. 拉取最新代码
git pull origin main

# 2. 根据更改选择更新方式：

# 仅前端更改
./scripts/deploy-docker.sh --frontend

# 仅后端更改  
./scripts/deploy-docker.sh --backend

# 前端+后端都有更改
./scripts/deploy-docker.sh

# 3. 验证部署
curl -f https://www.ifoodme.com/api/health
curl -f https://www.ifoodme.com/
```

### 完整重新部署
```bash
# 适用于重大更改或环境问题
./scripts/deploy-docker.sh
```

## 📊 性能优化

### 前端优化
- 静态文件直接通过 nginx 服务，性能更好
- 支持 gzip 压缩
- 静态资源缓存 1 年
- 支持现代浏览器的 HTTP/2

### 容器资源限制
```yaml
# 在 docker-compose.yml 中添加
deploy:
  resources:
    limits:
      cpus: '0.5'
      memory: 512M
```

## 🔧 脚本功能说明

### deploy-docker.sh 详细功能

**完整部署模式** (`./scripts/deploy-docker.sh`):
1. 环境检查 (Docker, Node.js, npm)
2. 前端构建 (npm install + npm run build)
3. 数据库备份 (如果存在)
4. 后端镜像构建
5. 服务启动 (postgres → ar-backend → nginx)
6. 健康检查

**前端更新模式** (`./scripts/deploy-docker.sh --frontend`):
1. 环境检查
2. 前端构建
3. nginx 重启 (如果运行中)

**后端更新模式** (`./scripts/deploy-docker.sh --backend`):
1. 环境检查
2. 后端镜像构建
3. 后端服务重启

## 📞 支持

### 开发工作流
1. **前端开发**: 在 `client/` 目录使用 `npm run dev`
2. **前端部署**: 运行 `./scripts/deploy-docker.sh --frontend`
3. **后端开发**: 修改代码后运行 `./scripts/deploy-docker.sh --backend`
4. **完整部署**: 运行 `./scripts/deploy-docker.sh`

### 日志位置
- nginx 访问日志: `nginx/logs/access.log`
- nginx 错误日志: `nginx/logs/error.log`
- 后端日志: `docker-compose logs ar-backend`
- 数据库日志: `docker-compose logs postgres`

### 快速命令参考
```bash
# 查看帮助
./scripts/deploy-docker.sh --help

# 完整部署
./scripts/deploy-docker.sh

# 快速前端更新
./scripts/deploy-docker.sh --frontend

# 后端更新
./scripts/deploy-docker.sh --backend

# 查看服务状态
docker-compose --env-file=.env.production ps

# 查看实时日志
docker-compose --env-file=.env.production logs -f
```

如果遇到问题，请按以下顺序检查：
1. 容器状态: `docker-compose ps`
2. 服务日志: `docker-compose logs <service_name>`
3. 网络连接: `curl` 测试各个端点
4. 文件权限: 确保静态文件可读取 