# 环境变量配置文档

本项目已完全迁移到使用环境变量进行配置管理，不再依赖硬编码的配置值。

## 📋 环境变量列表

### 🚀 服务器配置
| 变量名 | 描述 | 默认值 | 必需 |
|--------|------|--------|------|
| `SERVER_PORT` | 服务器端口 | `3000` | ❌ |
| `ENVIRONMENT` | 运行环境 (`development`/`production`) | `development` | ❌ |
| `FRONTEND_URL` | 前端应用URL | 根据环境自动设置 | ❌ |

### 🗄️ 数据库配置
| 变量名 | 描述 | 默认值 | 必需 |
|--------|------|--------|------|
| `BLUEPRINT_DB_HOST` | 数据库主机地址 | - | ✅ |
| `BLUEPRINT_DB_PORT` | 数据库端口 | - | ✅ |
| `BLUEPRINT_DB_USERNAME` | 数据库用户名 | - | ✅ |
| `BLUEPRINT_DB_PASSWORD` | 数据库密码 | - | ✅ |
| `BLUEPRINT_DB_DATABASE` | 数据库名称 | - | ✅ |
| `BLUEPRINT_DB_SCHEMA` | 数据库模式 | `public` | ❌ |

### 🔐 认证配置
| 变量名 | 描述 | 默认值 | 必需 |
|--------|------|--------|------|
| `JWT_SECRET` | JWT签名密钥 | - | ✅ |
| `JWT_REFRESH_SECRET` | JWT刷新令牌密钥 | `JWT_SECRET + "_refresh"` | ❌ |
| `SESSION_SECRET` | Session密钥 | `JWT_SECRET + "_session"` | ❌ |

### 🌐 OAuth配置
| 变量名 | 描述 | 默认值 | 必需 |
|--------|------|--------|------|
| `GOOGLE_CLIENT_ID` | Google OAuth客户端ID | - | ❌ |
| `GOOGLE_CLIENT_SECRET` | Google OAuth客户端密钥 | - | ❌ |
| `GOOGLE_CALLBACK_URL` | Google OAuth回调URL | 根据环境自动设置 | ❌ |

### 🔒 CORS和Cookie配置
| 变量名 | 描述 | 默认值 | 必需 |
|--------|------|--------|------|
| `ALLOWED_ORIGINS` | 允许的CORS域名（逗号分隔） | 根据环境自动设置 | ❌ |
| `COOKIE_DOMAIN` | Cookie域名 | 根据环境自动设置 | ❌ |

### 👤 管理员配置
| 变量名 | 描述 | 默认值 | 必需 |
|--------|------|--------|------|
| `ADMIN_EMAIL` | 管理员邮箱 | - | ❌ |
| `ADMIN_PASSWORD` | 管理员密码 | - | ❌ |

## 🔧 配置文件

### 开发环境 (`.env`)
```bash
# 服务器配置
SERVER_PORT=3000
ENVIRONMENT=development
FRONTEND_URL=http://localhost:3001

# 数据库配置
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_USERNAME=myuser
BLUEPRINT_DB_PASSWORD=mypassword
BLUEPRINT_DB_DATABASE=mydatabase
BLUEPRINT_DB_SCHEMA=public

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_CALLBACK_URL=http://localhost:3000/auth/google/callback

# JWT 密钥
JWT_SECRET=your_jwt_secret_key_for_development
JWT_REFRESH_SECRET=your_jwt_refresh_secret_key_for_development
SESSION_SECRET=your_session_secret_key_for_development

# CORS 配置
ALLOWED_ORIGINS=http://localhost:3001,http://localhost:3000

# Cookie 域名配置
COOKIE_DOMAIN=

# 管理员账户
ADMIN_EMAIL=admin@ar-backend.com
ADMIN_PASSWORD=admin123
```

### 生产环境示例
```bash
# 服务器配置
SERVER_PORT=3000
ENVIRONMENT=production
FRONTEND_URL=https://www.yourdomain.com

# 数据库配置
BLUEPRINT_DB_HOST=your-db-host
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_USERNAME=your_production_user
BLUEPRINT_DB_PASSWORD=your_secure_password
BLUEPRINT_DB_DATABASE=your_production_db

# Google OAuth
GOOGLE_CLIENT_ID=your_production_google_client_id
GOOGLE_CLIENT_SECRET=your_production_google_client_secret
GOOGLE_CALLBACK_URL=https://api.yourdomain.com/api/auth/google/callback

# JWT 密钥 (使用强密钥)
JWT_SECRET=your_very_secure_jwt_secret_key_here
JWT_REFRESH_SECRET=your_very_secure_refresh_secret_key_here
SESSION_SECRET=your_very_secure_session_secret_key_here

# CORS 配置
ALLOWED_ORIGINS=https://www.yourdomain.com,https://yourdomain.com

# Cookie 域名配置
COOKIE_DOMAIN=.yourdomain.com

# 管理员账户
ADMIN_EMAIL=admin@yourdomain.com
ADMIN_PASSWORD=your_secure_admin_password
```

## 🚨 安全注意事项

1. **密钥安全**: 所有密钥（JWT_SECRET、数据库密码等）应使用强随机字符串
2. **生产环境**: 生产环境中绝不要使用示例密钥
3. **环境隔离**: 开发、测试、生产环境应使用不同的密钥和配置
4. **版本控制**: `.env` 文件不应提交到版本控制系统
5. **权限管理**: 确保环境变量文件的访问权限受限

## 📝 配置验证

应用启动时会自动验证必需的环境变量：
- 如果缺少必需的环境变量，应用会报错并退出
- 可选的环境变量会使用合理的默认值
- 管理员账户会在首次启动时自动创建（如果配置了相关环境变量）

## 🔄 迁移说明

本次更新已完全移除：
- ❌ 硬编码的JWT密钥
- ❌ 硬编码的域名配置
- ❌ config.yaml中的数据库配置
- ❌ 硬编码的CORS域名
- ❌ 硬编码的Session密钥

所有配置现在都通过环境变量管理，提高了安全性和部署灵活性。 