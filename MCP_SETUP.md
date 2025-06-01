# MCP (Model Context Protocol) 设置指南

## 概述
本项目已配置MCP服务器来提供数据库和文件系统访问能力。

## 配置的MCP服务器

### 1. PostgreSQL 数据库服务器
```bash
npx -y @modelcontextprotocol/server-postgres postgresql://myuser:mypassword@localhost:5432/mydatabase
```

### 2. 文件系统服务器
```bash
npx -y @modelcontextprotocol/server-filesystem /Users/jamesw/dev_workspace/travel_AR_project_RXD/code/ar-backend
```

## 使用方法

### 启动MCP服务器

1. **启动数据库MCP服务器**：
```bash
npx -y @modelcontextprotocol/server-postgres postgresql://myuser:mypassword@localhost:5432/mydatabase
```

2. **启动文件系统MCP服务器**：
```bash
npx -y @modelcontextprotocol/server-filesystem /Users/jamesw/dev_workspace/travel_AR_project_RXD/code/ar-backend
```

### 数据库访问示例

查看所有表的记录数：
```sql
SELECT 
    'users' as table_name, COUNT(*) as row_count FROM users 
UNION ALL 
SELECT 'articles', COUNT(*) FROM articles 
-- ... 其他表
ORDER BY table_name;
```

查看用户数据：
```sql
SELECT user_id, name, email, provider, status, created_at FROM users LIMIT 10;
```

### 可用的数据库表

- `users` - 用户表
- `articles` - 文章表  
- `comments` - 评论表
- `facilities` - 设施表
- `files` - 文件表
- `languages` - 语言表
- `menus` - 菜单表
- `notices` - 通知表
- `refresh_tokens` - 刷新令牌表
- `stores` - 商店表
- `taggings` - 标签关联表
- `tags` - 标签表
- `visit_histories` - 访问历史表

## 安全规则

### 数据库访问限制
- ✅ 允许：SELECT, INSERT, UPDATE, DELETE
- ❌ 禁止：TRUNCATE, DROP
- ⚠️ 注意：批量操作前需备份

### 文件系统访问限制
- ✅ 允许的文件类型：`.go`, `.json`, `.yaml`, `.yml`, `.md`, `.txt`, `.sql`, `.sh`
- ❌ 排除的路径：`node_modules`, `.git`, `tmp`, `pgdata`, `pgadmin-data`, `*.log`, `ar-backend`

## 项目信息

- **名称**: AR Backend
- **描述**: Travel AR Project Backend API
- **技术栈**: Go, Gin, GORM, PostgreSQL, Docker
- **API地址**: http://localhost:3000/api
- **API文档**: http://localhost:3000/swagger/index.html
- **健康检查**: http://localhost:3000/api/health

## 开发命令

```bash
# 开发模式启动
air

# 构建项目
go build -o ar-backend .

# 部署
./deploy.sh

# 启动数据库
docker-compose up postgres -d

# 查看日志
tail -f app.log
```

## 故障排除

1. **数据库连接失败**：
   - 确保PostgreSQL容器正在运行：`docker ps | grep postgres`
   - 检查连接字符串中的用户名、密码和数据库名

2. **文件访问受限**：
   - 确保路径在允许的目录列表中
   - 检查文件扩展名是否在允许的列表中

3. **MCP服务器启动失败**：
   - 确保已安装Node.js和npm
   - 检查网络连接和防火墙设置 