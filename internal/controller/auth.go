package controller

import (
	"time"

	"ar-backend/internal/model"
	"ar-backend/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
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

// Register godoc
// @Summary 用户注册
// @Description 用户注册，创建新用户
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body model.UserReqCreate true "注册请求"
// @Success 200 {object} model.Response[model.RegisterResponse]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/auth/register [post]
func Register(c *gin.Context) {
	var req model.UserReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	db := database.GetDB()
	var count int64
	db.Model(&model.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "邮箱已被注册"})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "密码加密失败"})
		return
	}

	user := model.User{
		Name:        req.Name,
		NameKana:    req.NameKana,
		Address:     req.Address,
		Gender:      &req.Gender,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    string(hashedPwd),
		GoogleID:    req.GoogleID,
		AppleID:     req.AppleID,
		Provider:    req.Provider,
		Status:      req.Status,
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	token, err := generateTokenWithUserID(user.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "Token生成失败"})
		return
	}
	c.JSON(200, model.Response[model.RegisterResponse]{Success: true, Data: model.RegisterResponse{User: user, Token: token}})
}

func generateTokenWithUserID(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	type UserIDClaims struct {
		UserID int `json:"user_id"`
		jwt.RegisteredClaims
	}
	claims := &UserIDClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
