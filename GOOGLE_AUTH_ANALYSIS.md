# 🔐 Google 登录逻辑分析

## 📋 当前实现概述

项目实现了 **两套 Google 登录方式**，目前主要使用第一种：

### 🎯 方式一：Gothic OAuth 流程（主要使用）

#### 流程图
```
前端 → GET /api/auth/google → Google OAuth → /api/auth/google/callback → 后端处理 → 前端重定向
```

#### 详细步骤
1. **前端触发**: 用户点击 Google 登录按钮
   - 文件: `client/src/components/GoogleLoginButton.tsx`
   - 跳转: `window.location.href = "${apiBaseUrl}/api/auth/google"`

2. **后端处理开始认证**: 
   - 路由: `GET /api/auth/:provider` → `beginAuthProviderCallback()`
   - 文件: `internal/service/routes.go:215`
   - 使用: `gothic.BeginAuthHandler()` 重定向到 Google

3. **Google 认证**: 用户在 Google 页面完成认证

4. **回调处理**:
   - 路由: `GET /api/auth/:provider/callback` → `getAuthCallbackFunction()`
   - 文件: `internal/service/routes.go:114`
   - 使用: `gothic.CompleteUserAuth()` 获取用户信息

5. **用户处理**:
   - 检查用户是否存在 (`email` 字段)
   - 不存在则创建新用户，存在则更新信息
   - 生成 JWT Token (24小时有效期)
   - 设置 Cookie 和 URL 参数

6. **前端重定向**:
   - 重定向到 `FRONTEND_URL` 并携带 token
   - 前端接收 token 并存储到 localStorage

#### 配置文件
- **认证配置**: `internal/auth/auth.go`
- **路由配置**: `internal/service/routes.go`

### 🎯 方式二：ID Token 验证（备用）

#### 流程
```
前端获取ID Token → POST /api/auth/google → 后端验证 → 返回JWT
```

#### 实现
- 路由: `POST /api/auth/google` → `GoogleAuth()`
- 文件: `internal/controller/auth.go:373`
- 验证: 使用 Google API 验证 `id_token`

## 🔧 环境配置

### 开发环境 (.env)
```env

```

### 生产环境 (.env.production)
```env
GOOGLE_CLIENT_ID=your_production_google_client_id
GOOGLE_CLIENT_SECRET=your_production_google_client_secret
GOOGLE_CALLBACK_URL=https://www.ifoodme.com/api/auth/google/callback
COOKIE_DOMAIN=.ifoodme.com
```

## 🗄️ 数据库字段

### User 模型
```go
type User struct {
    UserID    int    `gorm:"primaryKey" json:"user_id"`
    Email     string `gorm:"unique;not null" json:"email"`
    Name      string `json:"name"`
    Avatar    string `json:"avatar"`
    GoogleID  string `json:"google_id"`
    Provider  string `json:"provider"`  // "google" | "email"
    Status    string `json:"status"`    // "active" | "pending"
    // ... 其他字段
}
```

## 🍪 Cookie 和 Token 处理

### JWT Token 字段
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "name": "User Name",
  "avatar": "avatar_url",
  "exp": 1748445335
}
```

### Cookie 设置
- **Name**: `token`
- **HttpOnly**: `false` (让前端 JS 可访问)
- **Secure**: 生产环境 `true`
- **SameSite**: 生产环境 `None`，开发环境 `Lax`
- **Domain**: 生产环境设置为 `.ifoodme.com`

## 🚨 已修复的问题

### ✅ 回调URL不匹配问题
**问题**: 开发环境默认回调URL缺少 `/api` 前缀
```go
// ❌ 错误 (之前)
callbackURL = "http://localhost:3000/auth/google/callback"

// ✅ 正确 (现在)
callbackURL = "http://localhost:3000/api/auth/google/callback"
```

**修复**: 更新了 `.env` 文件中的 `GOOGLE_CALLBACK_URL`

## 🔍 潜在优化点

### 1. 路由冲突
- 当前 `/api/auth/google` 既有 GET 又有 POST
- 建议将 ID Token 方式改为 `/api/auth/google/token`

### 2. Token 一致性
- Gothic 流程使用简单 JWT (24h)
- ID Token 流程使用 Access + Refresh Token
- 建议统一使用 Refresh Token 机制

### 3. 错误处理
- 增加更详细的错误日志
- 添加认证失败的前端友好提示

## 🧪 测试建议

### 开发环境测试
1. 启动后端: `air` 或 `go run main.go`
2. 访问前端登录页面
3. 点击 Google 登录按钮
4. 检查控制台是否有 token
5. 验证 `/api/me` 接口返回用户信息

### 生产环境注意事项
1. 确保 Google Cloud Console 中配置了正确的回调URL
2. 验证 HTTPS 和 Cookie 安全设置
3. 测试跨域访问是否正常 

graph TD
    A[用户点击Google登录] --> B[跳转到 https://www.ifoodme.com/api/auth/google]
    B --> C[后端beginAuthProviderCallback调用gothic.BeginAuthHandler]
    C --> D[重定向到Google OAuth页面]
    D --> E[用户在Google完成认证]
    E --> F[Google回调到 /api/auth/google/callback]
    F --> G[后端getAuthCallbackFunction调用gothic.CompleteUserAuth]
    G --> H[获取Google用户信息]
    H --> I[检查用户是否存在数据库中]
    I --> J[创建新用户或更新现有用户]
    J --> K[生成JWT Token 24小时有效期]
    K --> L[设置Cookie并添加URL参数]
    L --> M[重定向到FRONTEND_URL/?token=xxx]
    M --> N[前端从URL获取token并存储到localStorage] 