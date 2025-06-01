# 🔄 Redirect 参数修改说明

## 📋 修改概述

修改了后端 Google OAuth 流程，现在**优先使用前端传递的 `redirect` 参数**进行登录后的跳转。

## 🔄 修改详情

### 1. **beginAuthProviderCallback 函数修改**
- **位置**: `internal/service/routes.go`
- **功能**: 获取前端传递的 `redirect` 参数并保存到 session

```go
// 获取前端传递的 redirect 参数
redirectURL := c.Query("redirect")
if redirectURL != "" {
    // 将 redirect URL 保存到 session 中，以便在回调时使用
    session, err := gothic.Store.Get(r, "oauth_session")
    if err == nil {
        session.Values["redirect_url"] = redirectURL
        session.Save(r, w)
    }
}
```

### 2. **getAuthCallbackFunction 函数修改**
- **位置**: `internal/service/routes.go`
- **功能**: 按优先级获取重定向URL

#### 优先级顺序
1. **🥇 Session 中保存的 redirect URL** (前端传递)
2. **🥈 环境变量 FRONTEND_URL**
3. **🥉 默认前端URL**

```go
var frontendURL string

// 1. 优先从 session 中获取前端传递的 redirect 参数
session, err := gothic.Store.Get(r, "oauth_session")
if err == nil {
    if savedRedirectURL, ok := session.Values["redirect_url"].(string); ok && savedRedirectURL != "" {
        frontendURL = savedRedirectURL
    }
}

// 2. 如果没有 redirect 参数，使用环境变量
if frontendURL == "" {
    frontendURL = os.Getenv("FRONTEND_URL")
}

// 3. 最后使用默认地址
if frontendURL == "" {
    frontendURL = getDefaultFrontendURL()
}
```

## 🔗 前端使用方式

### 当前前端代码
```typescript
const handleLogin = () => {
    const currentURL = window.location.origin;
    const redirectParam = encodeURIComponent(currentURL);
    const apiBaseUrl = 'https://www.ifoodme.com';
    
    // redirect 参数现在会被后端使用
    window.location.href = `${apiBaseUrl}/api/auth/google?redirect=${redirectParam}`;
};
```

### 支持的 redirect 格式
- **完整URL**: `https://example.com/dashboard`
- **相对路径**: `/dashboard`
- **带参数**: `https://example.com/page?param=value`

## 🧪 测试场景

### 场景1: 前端传递 redirect 参数
```
输入: /api/auth/google?redirect=https://example.com/dashboard
结果: 登录后跳转到 https://example.com/dashboard?token=xxx
```

### 场景2: 没有 redirect 参数
```
输入: /api/auth/google
结果: 使用环境变量 FRONTEND_URL 或默认URL
```

### 场景3: redirect 参数为空
```
输入: /api/auth/google?redirect=
结果: 使用环境变量 FRONTEND_URL 或默认URL
```

## 🛡️ 安全考虑

### 1. **URL 验证**
当前实现直接使用前端传递的 redirect URL，建议添加验证：

```go
// 推荐添加的安全验证
func isValidRedirectURL(redirectURL string) bool {
    allowedDomains := []string{
        "localhost",
        "ifoodme.com",
        "www.ifoodme.com",
    }
    
    u, err := url.Parse(redirectURL)
    if err != nil {
        return false
    }
    
    for _, domain := range allowedDomains {
        if strings.Contains(u.Host, domain) {
            return true
        }
    }
    
    return false
}
```

### 2. **防止开放重定向攻击**
- 验证重定向URL的域名
- 只允许白名单中的域名
- 拒绝外部不可信域名

## 📊 日志输出

修改后会输出更详细的日志：

```
Begin Auth - Provider: google
前端传递的redirect参数: https://example.com/dashboard
成功保存redirect_url到session: https://example.com/dashboard

...（OAuth流程）...

OAuth Callback - Provider: google
使用前端传递的redirect URL: https://example.com/dashboard
最终重定向到: https://example.com/dashboard?token=eyJhbGc...
```

## 🚀 部署注意事项

1. **Session Store**: 确保 Gothic Session Store 正确配置
2. **Cookie 设置**: 检查跨域 Cookie 设置
3. **HTTPS**: 生产环境确保使用 HTTPS
4. **域名验证**: 建议添加 redirect URL 白名单验证

## 🔧 后续优化建议

1. **添加 URL 验证**: 防止开放重定向攻击
2. **Session 清理**: 成功重定向后清理 session 数据
3. **错误处理**: 改善 session 操作的错误处理
4. **日志级别**: 生产环境可考虑降低日志详细程度 