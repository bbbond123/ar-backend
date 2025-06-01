# AR Backend 部署指南

## 🔒 安全部署流程

### 1. 环境配置

首先复制环境变量示例文件：
```bash
cp .env.example .env
```

编辑`.env`文件，填入实际的配置信息：
```bash
# 服务器配置
SERVER_PORT=3000
ENVIRONMENT=production

# 数据库配置 (替换为实际值)
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_USERNAME=your_actual_username
BLUEPRINT_DB_PASSWORD=your_secure_password
BLUEPRINT_DB_DATABASE=your_database_name

# JWT密钥 (使用强密码)
JWT_SECRET=your_very_secure_jwt_secret

# 管理员账户
ADMIN_EMAIL=your_admin@example.com
ADMIN_PASSWORD=your_secure_admin_password
```

### 2. 数据库启动

启动PostgreSQL容器：
```bash
docker-compose up postgres -d
```

### 3. 部署选项

#### 首次部署 (包含数据初始化)
```bash
./deploy.sh --init-db
```

#### 后续部署 (保留现有数据)
```bash
./deploy.sh
```

### 4. 验证部署

检查服务状态：
```bash
curl http://localhost:3000/api/health
```

查看API文档：
```bash
open http://localhost:3000/swagger/index.html
```

## 📊 数据管理

### 初始化数据库
如果需要重新初始化数据库数据：
```bash
./scripts/init_db.sh
```

### 备份数据库
```bash
docker exec ar-backend-postgres pg_dump -U $BLUEPRINT_DB_USERNAME $BLUEPRINT_DB_DATABASE > backup.sql
```

### 恢复数据库
```bash
docker exec -i ar-backend-postgres psql -U $BLUEPRINT_DB_USERNAME $BLUEPRINT_DB_DATABASE < backup.sql
```

## 🔐 安全注意事项

### 环境变量
- ✅ `.env`文件已在`.gitignore`中，不会被提交
- ✅ 敏感信息通过环境变量配置
- ✅ 生产环境使用强密码

### 数据库
- ✅ 数据库初始化脚本确保基础数据存在
- ✅ 管理员账户使用bcrypt加密
- ✅ 支持数据备份和恢复

### 文件系统
- ✅ Docker数据卷持久化
- ✅ 日志文件被忽略，不会提交敏感信息

## 📁 文件结构

```
ar-backend/
├── .env.example          # 环境变量示例
├── .env                  # 实际环境变量 (不提交)
├── .gitignore           # 忽略敏感文件
├── deploy.sh            # 部署脚本
├── scripts/
│   ├── init_db.sh      # 数据库初始化脚本
│   └── init_data.sql   # 初始化SQL
├── .mcp.json           # MCP配置 (使用环境变量)
└── DEPLOYMENT.md       # 本文档
```

## 🚀 常用命令

```bash
# 开发模式
air

# 构建项目
go build -o ar-backend .

# 首次部署
./deploy.sh --init-db

# 常规部署
./deploy.sh

# 查看日志
tail -f app.log

# 重启数据库
docker-compose restart postgres

# 数据库连接
docker exec -it ar-backend-postgres psql -U $BLUEPRINT_DB_USERNAME $BLUEPRINT_DB_DATABASE
```

## 🛠️ 故障排除

### 服务启动失败
1. 检查端口是否被占用：`lsof -i :3000`
2. 查看应用日志：`tail -f app.log`
3. 检查环境变量：`cat .env`

### 数据库连接失败
1. 确保容器运行：`docker ps | grep postgres`
2. 检查数据库配置：`.env`文件中的数据库设置
3. 测试连接：`docker exec ar-backend-postgres pg_isready`

### 数据丢失
1. 检查Docker卷：`docker volume ls | grep postgres`
2. 恢复备份：使用上面的恢复命令
3. 重新初始化：`./scripts/init_db.sh` 