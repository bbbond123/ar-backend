package server

import (
	"ar-backend/internal/model"
	"ar-backend/internal/router"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/markbates/goth/gothic"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// getJWTSecret 从环境变量获取 JWT 密钥
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	return []byte(secret)
}

// getAllowedOrigins 从环境变量获取允许的 CORS 域名
func getAllowedOrigins() []string {
	// 从环境变量获取允许的域名，用逗号分隔
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	if originsEnv != "" {
		return strings.Split(originsEnv, ",")
	}

	// 如果没有设置环境变量，使用默认值
	if os.Getenv("ENVIRONMENT") == "production" {
		return []string{
			os.Getenv("FRONTEND_URL"),
			"https://ifoodme.com",
			"https://www.ifoodme.com",
		}
	}

	// 开发环境默认值
	return []string{
		"http://localhost:3001",
		"http://localhost:3000",
		os.Getenv("FRONTEND_URL"),
	}
}

// getDomainFromEnv 从环境变量获取域名
func getDomainFromEnv() string {
	domain := os.Getenv("COOKIE_DOMAIN")
	if domain == "" && os.Getenv("ENVIRONMENT") == "production" {
		domain = ".ifoodme.com"
	}
	return domain
}

// getDefaultFrontendURL 获取默认前端 URL
func getDefaultFrontendURL() string {
	if os.Getenv("ENVIRONMENT") == "production" {
		return "https://www.ifoodme.com/"
	}
	return "http://localhost:3001/"
}

func (s *Server) RegisterRoutes() http.Handler {
	// r := gin.Default()
	r := router.InitRouter()

	// 动态获取允许的域名
	allowedOrigins := getAllowedOrigins()

	// CORS 配置
	corsConfig := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "x-app-platform"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	r.GET("/api", s.HelloWorldHandler)
	r.GET("/api/health", s.healthHandler)

	r.GET("/api/auth/:provider", s.beginAuthProviderCallback)
	r.GET("/api/auth/:provider/callback", s.getAuthCallbackFunction)
	r.GET("/api/users/me", s.MeHandler)

	// r.POST("/api/logout", s.LogoutHandler)
	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) getAuthCallbackFunction(c *gin.Context) {
	fmt.Printf("\n=== GGG getAuthCallbackFunction 开始 ===\n")
	provider := c.Param("provider")
	ctx := context.WithValue(context.Background(), "provider", provider)
	r := c.Request.WithContext(ctx)
	w := c.Writer

	fmt.Printf("OAuth Callback - Provider: %s\n", provider)
	fmt.Printf("OAuth Callback - Request Host: %s\n", r.Host)
	fmt.Printf("OAuth Callback - Request URL: %s\n", r.URL.String())
	fmt.Printf("OAuth Callback - Request Cookies: %v\n", r.Cookies())
	fmt.Printf("OAuth Callback - State: %s\n", r.URL.Query().Get("state"))
	fmt.Printf("OAuth Callback - Code: %s\n", r.URL.Query().Get("code"))
	fmt.Printf("OAuth Callback - Error: %s\n", r.URL.Query().Get("error"))

	fmt.Printf("开始完成Gothic OAuth认证...\n")
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Printf("❌ OAuth Callback Error: %v\n", err)
		c.String(http.StatusUnauthorized, "auth error: %v", err)
		return
	}

	fmt.Printf("✅ Gothic OAuth认证成功\n")
	fmt.Printf("获取到的用户信息 - Email: %s, Name: %s, UserID: %s\n", user.Email, user.Name, user.UserID)
	fmt.Printf("Avatar: %s, Provider: %s\n", user.AvatarURL, user.Provider)

	var userInDB model.User
	err = s.gormDB.Where("email = ?", user.Email).First(&userInDB).Error
	if err == gorm.ErrRecordNotFound {
		fmt.Printf("🆕 用户不存在，创建新用户 - Email: %s\n", user.Email)
		userInDB = model.User{
			Email:    user.Email,
			GoogleID: user.UserID,
			Name:     user.Name,
			Avatar:   user.AvatarURL,
			Provider: "google",
			Status:   "active",
		}
		err = s.gormDB.Create(&userInDB).Error
		if err != nil {
			fmt.Printf("❌ 创建用户失败: %v\n", err)
			c.String(http.StatusInternalServerError, "Could not create user")
			return
		}
		fmt.Printf("✅ 新用户创建成功 - UserID: %d\n", userInDB.UserID)
	} else if err == nil {
		fmt.Printf("✅ 找到已存在用户 - UserID: %d, Email: %s\n", userInDB.UserID, userInDB.Email)
		err = s.gormDB.Model(&userInDB).Updates(map[string]interface{}{
			"google_id":  user.UserID,
			"avatar":     user.AvatarURL,
			"updated_at": time.Now(),
		}).Error
		if err != nil {
			fmt.Printf("❌ 更新用户信息失败: %v\n", err)
			c.String(http.StatusInternalServerError, "Could not update user")
			return
		}
		fmt.Printf("✅ 用户信息更新成功\n")
	} else {
		fmt.Printf("❌ 数据库查询错误: %v\n", err)
		c.String(http.StatusInternalServerError, "Database error")
		return
	}

	fmt.Printf("开始生成JWT Token - UserID: %d\n", userInDB.UserID)

	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userInDB.UserID,
		"email":   userInDB.Email,
		"name":    userInDB.Name,
		"avatar":  userInDB.Avatar,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		fmt.Printf("❌ JWT Token生成失败: %v\n", err)
		c.String(http.StatusInternalServerError, "Could not create token")
		return
	}

	fmt.Printf("✅ JWT Token生成成功，长度: %d\n", len(tokenString))

	// 设置 Cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: false,                 // 改为 false，让前端JS能访问
		Secure:   true,                  // 生产环境用 true（HTTPS）
		SameSite: http.SameSiteNoneMode, // 跨域需要 None
	}

	// 检查环境并设置Cookie安全选项
	isProd := os.Getenv("ENVIRONMENT") == "production"
	if isProd {
		cookieDomain := os.Getenv("COOKIE_DOMAIN")
		if cookieDomain != "" {
			cookie.Domain = cookieDomain
		}
		fmt.Printf("设置生产环境Cookie - Domain: %s, Secure: true\n", cookie.Domain)
	} else {
		cookie.Secure = false
		fmt.Printf("设置开发环境Cookie - Secure: false\n")
	}

	// 不设置域名，让浏览器自动处理
	http.SetCookie(c.Writer, cookie)
	fmt.Printf("✅ Cookie设置成功\n")

	// 获取重定向地址
	var frontendURL string

	// 1. 优先从 session 中获取前端传递的 redirect 参数
	session, err := gothic.Store.Get(r, "oauth_session")
	if err == nil {
		fmt.Printf("✅ 成功获取session\n")
		fmt.Printf("Session所有值: %+v\n", session.Values)
		if savedRedirectURL, ok := session.Values["redirect_url"].(string); ok && savedRedirectURL != "" {
			frontendURL = savedRedirectURL
			fmt.Printf("✅ 从session获取到前端传递的redirect URL: %s\n", frontendURL)
		} else {
			fmt.Printf("❌ session中没有找到redirect_url或值为空\n")
			fmt.Printf("redirect_url值类型: %T, 值: %v\n", session.Values["redirect_url"], session.Values["redirect_url"])
		}
	} else {
		fmt.Printf("❌ 获取session失败: %v\n", err)
	}

	// 2. 如果没有 redirect 参数，优先从环境变量获取
	if frontendURL == "" {
		frontendURL = os.Getenv("FRONTEND_URL")
		if frontendURL != "" {
			fmt.Printf("使用环境变量FRONTEND_URL: %s\n", frontendURL)
		}
	}

	// 3. 如果环境变量也没有，使用默认地址
	if frontendURL == "" {
		frontendURL = getDefaultFrontendURL()
		fmt.Printf("使用默认前端URL: %s\n", frontendURL)
	}

	// 检查是否是React Native应用的深度链接
	if strings.HasPrefix(frontendURL, "travelview://") {
		fmt.Printf("🔗 AAAA 检测到React Native深度链接\n")
		// React Native 深度链接，构造参数
		deepLink := frontendURL + "?token=" + url.QueryEscape(tokenString) +
			"&user_id=" + fmt.Sprintf("%d", userInDB.UserID) +
			"&email=" + url.QueryEscape(userInDB.Email) +
			"&name=" + url.QueryEscape(userInDB.Name) +
			"&code=" + url.QueryEscape(r.URL.Query().Get("code")) +
			"&state=" + url.QueryEscape(r.URL.Query().Get("state"))

		if userInDB.Avatar != "" {
			deepLink += "&avatar=" + url.QueryEscape(userInDB.Avatar)
		}

		fmt.Printf("构造的深度链接: %s\n", deepLink)
		fmt.Printf("=== RRR 重定向到React Native App ===\n\n")

		// 重定向到深度链接
		c.Redirect(http.StatusFound, deepLink)
		return
	}

	// Web应用处理
	fmt.Printf("🌐 处理Web应用重定向\n")
	// 确保URL以斜杠结尾
	if !strings.HasSuffix(frontendURL, "/") {
		frontendURL += "/"
	}

	// 同时将token作为URL参数传递，让前端可以获取并存储
	frontendURL += "?token=" + url.QueryEscape(tokenString)

	fmt.Printf("最终重定向到: %s\n", frontendURL)
	fmt.Printf("=== getAuthCallbackFunction 成功完成 ===\n")
	fmt.Printf("用户: %s (ID: %d) 通过OAuth登录成功\n\n", userInDB.Email, userInDB.UserID)

	// 重定向到前端首页
	c.Redirect(http.StatusFound, frontendURL)
}

func (s *Server) beginAuthProviderCallback(c *gin.Context) {
	fmt.Printf("\n=== X beginAuthProviderCallback 开始 ===\n")
	provider := c.Param("provider")
	ctx := context.WithValue(context.Background(), "provider", provider)
	r := c.Request.WithContext(ctx)
	w := c.Writer

	fmt.Printf("Begin Auth - Provider: %s\n", provider)
	fmt.Printf("Begin Auth - Request Host: %s\n", r.Host)
	fmt.Printf("Begin Auth - Request URL: %s\n", r.URL.String())
	fmt.Printf("Begin Auth - Request Cookies: %v\n", r.Cookies())
	fmt.Printf("Begin Auth - User-Agent: %s\n", r.Header.Get("User-Agent"))
	fmt.Printf("Begin Auth - Referer: %s\n", r.Header.Get("Referer"))

	// 获取前端传递的 redirect 参数
	// 支持两种参数名：redirect 和 redirect_uri
	redirectURL := c.Query("redirect")
	if redirectURL == "" {
		redirectURL = c.Query("redirect_uri")
	}
	fmt.Printf("接收到的redirect参数: %s\n", redirectURL)
	fmt.Printf("2222原始查询字符串: %s\n", c.Request.URL.RawQuery)

	if redirectURL != "" {
		fmt.Printf("处理前端传递的redirect参数: %s\n", redirectURL)

		// 将 redirect URL 保存到 session 中，以便在回调时使用
		session, err := gothic.Store.Get(r, "oauth_session")
		if err != nil {
			fmt.Printf("❌ 无法获取session: %v\n", err)
		} else {
			session.Values["redirect_url"] = redirectURL
			err = session.Save(r, w)
			if err != nil {
				fmt.Printf("❌ 无法保存session: %v\n", err)
			} else {
				fmt.Printf("✅ 成功保存redirect_url到session: %s\n", redirectURL)
			}
		}
	} else {
		fmt.Printf("无redirect参数，将使用默认重定向\n")
	}

	fmt.Printf("准备重定向到Google OAuth页面...\n")
	fmt.Printf("=== beginAuthProviderCallback 结束 ===\n\n")

	gothic.BeginAuthHandler(w, r)
}

func (s *Server) MeHandler(c *gin.Context) {
	var tokenStr string

	// 优先从 Authorization Header 读取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Printf("从Authorization Header获取token: %s...\n", tokenStr[:20])
	} else {
		// 从 Cookie 读取
		cookie, err := c.Request.Cookie("token")
		if err != nil {
			fmt.Printf("无法获取token: Authorization Header为空且Cookie读取失败: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		tokenStr = cookie.Value
		fmt.Printf("从Cookie获取token: %s...\n", tokenStr[:20])
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		fmt.Printf("JWT验证，使用secret: %s\n", string(getJWTSecret()))
		return getJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		fmt.Printf("JWT验证失败: err=%v, valid=%v\n", err, token.Valid)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Printf("JWT claims转换失败\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		fmt.Printf("无法获取user_id from claims: %v\n", claims)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	fmt.Printf("JWT验证成功，user_id: %v\n", userID)

	var user model.User
	err = s.gormDB.Where("user_id = ?", int(userID)).First(&user).Error
	if err != nil {
		fmt.Printf("数据库查询用户失败: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	fmt.Printf("用户查询成功: %s\n", user.Email)

	// c.JSON(200, model.Response[model.AuthResponse]{
	// 	Success: true,
	// 	Code:    200,
	// 	Data: model.AuthResponse{
	// 		User:         user,
	// 		AccessToken:  accessToken,
	// 		RefreshToken: refreshToken,
	// 	},
	// })

	c.JSON(http.StatusOK, model.Response[model.User]{
		Success: true,
		Code:    200,
		Data:    user,
	})
}

func (s *Server) LogoutHandler(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true, // 生产环境要 true
		SameSite: http.SameSiteNoneMode,
	}

	// 从环境变量设置域名
	domain := getDomainFromEnv()
	if domain != "" {
		cookie.Domain = domain
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
