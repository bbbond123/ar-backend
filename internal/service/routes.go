package server

import (
	"ar-backend/internal/model"
	"ar-backend/internal/router"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
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

// 假设你有一个 JWT secret
var jwtSecret = []byte("my_secret_key")

func (s *Server) RegisterRoutes() http.Handler {
	// r := gin.Default()
	r := router.InitRouter()

	// CORS 配置
	corsConfig := cors.Config{
		AllowOrigins: []string{
			"https://ifoodme.com",
			"https://www.ifoodme.com",
			"https://api.ifoodme.com",
			"http://localhost:5173", // 允许本地前端
			"http://localhost:3000", // 允许本地后端
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	r.GET("/api", s.HelloWorldHandler)
	r.GET("/api/health", s.healthHandler)

	r.GET("/api/auth/:provider", s.beginAuthProviderCallback)
	r.GET("/api/auth/:provider/callback", s.getAuthCallbackFunction)
	r.GET("/api/me", s.MeHandler)

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
	provider := c.Param("provider")
	ctx := context.WithValue(context.Background(), "provider", provider)
	r := c.Request.WithContext(ctx)
	w := c.Writer

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		c.String(http.StatusUnauthorized, "auth error: %v", err)
		return
	}
	var userInDB model.User
	err = s.gormDB.Where("email = ?", user.Email).First(&userInDB).Error
	if err == gorm.ErrRecordNotFound {
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
			c.String(http.StatusInternalServerError, "Could not create user")
			return
		}
	} else if err == nil {
		err = s.gormDB.Model(&userInDB).Updates(map[string]interface{}{
			"google_id":  user.UserID,
			"avatar":     user.AvatarURL,
			"updated_at": time.Now(),
		}).Error
		if err != nil {
			c.String(http.StatusInternalServerError, "Could not update user")
			return
		}
	} else {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	// 生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userInDB.UserID,
		"email":   userInDB.Email,
		"name":    userInDB.Name,
		"avatar":  userInDB.Avatar,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not create token")
		return
	}

	// 设置 Cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: false,                 // 改为 false，让前端JS能访问
		Secure:   true,                  // 生产环境用 true（HTTPS）
		SameSite: http.SameSiteNoneMode, // 跨域需要 None
	}

	// 不设置域名，让浏览器自动处理
	http.SetCookie(c.Writer, cookie)

	// 动态获取前端重定向地址
	var frontendURL string

	// 1. 优先从 Referer 头获取（前端发起登录请求的地址）
	referer := c.GetHeader("Referer")
	if referer != "" {
		// 安全检查：只允许特定的域名
		allowedDomains := []string{
			"http://localhost:5173",
			"http://localhost:3000",
			"https://ifoodme.com",
			"https://www.ifoodme.com",
		}

		for _, domain := range allowedDomains {
			if strings.HasPrefix(referer, domain) {
				frontendURL = referer
				break
			}
		}
	}

	// 2. 如果 Referer 不可用，使用环境变量
	if frontendURL == "" {
		frontendURL = os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			// 3. 默认地址，根据环境判断
			if os.Getenv("ENVIRONMENT") == "production" {
				frontendURL = "https://ifoodme.com/"
			} else {
				frontendURL = "http://localhost:5173/"
			}
		}
	}

	// 同时将token作为URL参数传递，让前端可以获取并存储
	frontendURL += "?token=" + url.QueryEscape(tokenString)

	fmt.Printf("重定向到: %s\n", frontendURL)

	// 重定向到前端首页
	c.Redirect(http.StatusFound, frontendURL)
}

func (s *Server) beginAuthProviderCallback(c *gin.Context) {
	provider := c.Param("provider")
	ctx := context.WithValue(context.Background(), "provider", provider)
	r := c.Request.WithContext(ctx)
	w := c.Writer
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
		fmt.Printf("JWT验证，使用secret: %s\n", string(jwtSecret))
		return jwtSecret, nil
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

	c.JSON(http.StatusOK, gin.H{
		"user_id":  user.UserID,
		"email":    user.Email,
		"name":     user.Name,
		"avatar":   user.Avatar,
		"provider": user.Provider,
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

	// 生产环境设置域名
	if c.Request.Host == "api.ifoodme.com" {
		cookie.Domain = ".ifoodme.com"
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
