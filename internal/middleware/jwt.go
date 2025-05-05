package middleware

import (
	"ar-backend/internal/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const secretKey = "my_secret_key" // 保持和主项目一致

type UserIDClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.BaseResponse{Success: false, ErrMessage: "未登录，缺少token"})
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &UserIDClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, model.BaseResponse{Success: false, ErrMessage: "token无效或已过期"})
			c.Abort()
			return
		}
		// 用户ID写入上下文
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
