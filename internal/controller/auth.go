package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Login godoc
// @Summary 登录
// @Description 登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body model.LoginRequest true "登录请求"
// @Success 200 {object} model.Response[model.LoginResponse]
// @Failure 401 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username != "admin" || password != "password" {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}
	token, _ := generateToken(username)
	c.JSON(200, gin.H{"token": token})
}

func RefreshToken(c *gin.Context) {
	tokenStr := c.PostForm("token")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "invalid token"})
		return
	}

	newToken, _ := generateToken(claims.Username)
	c.JSON(200, gin.H{"token": newToken})
}
