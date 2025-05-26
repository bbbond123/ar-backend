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
var jwtSecret = []byte("your_secret_key")

func (s *Server) RegisterRoutes() http.Handler {
	// r := gin.Default()
	r := router.InitRouter()

	// CORS 配置
	corsConfig := cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
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

	// 获取深层链接
	redirectURL := "com.travelview.app"
	if redirectURL == "" {
		redirectURL = "com.travelview.app:/oauth2redirect/google" // 默认深层链接
	}

	// 通过查询参数返回 token
	redirectURLWithToken := fmt.Sprintf("%s?token=%s", redirectURL, url.QueryEscape(tokenString))

	// 返回 JSON 响应（推荐）
	c.JSON(http.StatusOK, gin.H{
		"token":       tokenString,
		"redirectURL": redirectURLWithToken,
		"user": gin.H{
			"email":   userInDB.Email,
			"name":    userInDB.Name,
			"avatar":  userInDB.Avatar,
			"user_id": userInDB.UserID,
		},
	})

	// 前端返回 token 和 redirectURL
	// // 设置 Cookie
	// http.SetCookie(c.Writer, &http.Cookie{
	// 	Name:     "token",
	// 	Value:    tokenString,
	// 	Path:     "/",
	// 	HttpOnly: true,
	// 	Secure:   false, // 本地开发用 false，生产环境要 true
	// 	SameSite: http.SameSiteLaxMode,
	// })
	// // 重定向到前端
	// c.Redirect(http.StatusFound, "https://api.ifoodme.com/")
}

func (s *Server) beginAuthProviderCallback(c *gin.Context) {
	provider := c.Param("provider")
	ctx := context.WithValue(context.Background(), "provider", provider)
	r := c.Request.WithContext(ctx)
	w := c.Writer
	gothic.BeginAuthHandler(w, r)
}

func (s *Server) MeHandler(c *gin.Context) {
	cookie, err := c.Request.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	tokenStr := cookie.Value
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var user model.User
	err = s.gormDB.Where("user_id = ?", int(userID)).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":  user.UserID,
		"email":    user.Email,
		"name":     user.Name,
		"avatar":   user.Avatar,
		"provider": user.Provider,
	})
}

func (s *Server) LogoutHandler(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // 生产环境要 true
		SameSite: http.SameSiteLaxMode,
	})
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
