package middleware

import (
	"ar-backend/internal/model"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// getJWTSecret 从环境变量获取 JWT 密钥
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	return []byte(secret)
}

type UserIDClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("\n=== JWT中间件验证开始 ===\n")
		fmt.Printf("请求路径: %s %s\n", c.Request.Method, c.Request.URL.Path)
		fmt.Printf("请求IP: %s\n", c.ClientIP())
		fmt.Printf("User-Agent: %s\n", c.GetHeader("User-Agent"))
		
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("Authorization Header: %s\n", func() string {
			if authHeader == "" {
				return "空"
			}
			if len(authHeader) > 50 {
				return authHeader[:50] + "..."
			}
			return authHeader
		}())
		
		if authHeader == "" {
			fmt.Printf("❌ JWT验证失败: 缺少Authorization Header\n")
			c.JSON(http.StatusUnauthorized, model.BaseResponse{Success: false, ErrMessage: "未登录，缺少token"})
			c.Abort()
			return
		}
		
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Printf("提取的Token长度: %d\n", len(tokenStr))
		fmt.Printf("Token前30字符: %s...\n", func() string {
			if len(tokenStr) > 30 {
				return tokenStr[:30]
			}
			return tokenStr
		}())
		
		claims := &UserIDClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			fmt.Printf("JWT解析 - 使用算法: %v\n", token.Header["alg"])
			return getJWTSecret(), nil
		})
		
		if err != nil {
			fmt.Printf("❌ JWT解析失败: %v\n", err)
			c.JSON(http.StatusUnauthorized, model.BaseResponse{Success: false, ErrMessage: "token解析失败: " + err.Error()})
			c.Abort()
			return
		}
		
		if !token.Valid {
			fmt.Printf("❌ JWT验证失败: token无效\n")
			c.JSON(http.StatusUnauthorized, model.BaseResponse{Success: false, ErrMessage: "token无效"})
			c.Abort()
			return
		}
		
		fmt.Printf("✅ JWT验证成功\n")
		fmt.Printf("解析的Claims - UserID: %d\n", claims.UserID)
		fmt.Printf("Token过期时间: %v\n", claims.ExpiresAt)
		
		// 用户ID写入上下文
		c.Set("user_id", claims.UserID)
		fmt.Printf("✅ 用户ID已写入上下文: %d\n", claims.UserID)
		fmt.Printf("=== JWT中间件验证完成 ===\n\n")
		
		c.Next()
	}
}
