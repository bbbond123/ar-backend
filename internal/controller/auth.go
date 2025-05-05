package controller

import (
	"time"

	"ar-backend/internal/model"
	"ar-backend/pkg/database"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Login godoc
// @Summary 登录
// @Description 登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body model.LoginRequest true "登录请求"
// @Success 200 {object} model.Response[model.AuthResponse]
// @Failure 400 {object} model.BaseResponse
// @Failure 401 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var user model.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, model.BaseResponse{Success: false, ErrMessage: "用户不存在"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(401, model.BaseResponse{Success: false, ErrMessage: "密码错误"})
		return
	}
	accessToken, err := generateTokenWithUserID(user.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "Token生成失败"})
		return
	}
	refreshToken := generateRefreshToken()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := db.Create(&model.RefreshToken{
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Revoked:      false,
	}).Error; err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "RefreshToken存储失败"})
		return
	}
	c.JSON(200, model.Response[model.AuthResponse]{
		Success: true,
		Data: model.AuthResponse{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册，创建新用户
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body model.RegisterRequest true "注册请求"
// @Success 200 {object} model.Response[model.AuthResponse]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/auth/register [post]
func Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "参数错误: " + err.Error()})
		return
	}

	// 邮箱格式校验
	if !govalidator.IsEmail(req.Email) {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "邮箱格式不正确"})
		return
	}

	// 密码强度校验
	if len(req.Password) < 6 {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "密码长度不能少于6位"})
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
		Email:    req.Email,
		Password: string(hashedPwd),
		Provider: "email",
		Status:   "1",
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	accessToken, err := generateTokenWithUserID(user.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "Token生成失败"})
		return
	}

	refreshToken := generateRefreshToken()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := db.Create(&model.RefreshToken{
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Revoked:      false,
	}).Error; err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "RefreshToken存储失败"})
		return
	}

	c.JSON(200, model.Response[model.AuthResponse]{
		Success: true,
		Data: model.AuthResponse{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
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

// 控制器
func RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var token model.RefreshToken
	if err := db.Where("refresh_token = ? AND revoked = false AND expires_at > ?", req.RefreshToken, time.Now()).First(&token).Error; err != nil {
		c.JSON(401, model.BaseResponse{Success: false, ErrMessage: "refresh token无效"})
		return
	}
	// 生成新access token
	accessToken, err := generateTokenWithUserID(token.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "Token生成失败"})
		return
	}
	c.JSON(200, model.Response[model.RefreshTokenResponse]{Success: true, Data: model.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // 或生成新refresh token
	}})
}

func RevokeRefreshToken(c *gin.Context) {
	var req model.RevokeTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	db.Model(&model.RefreshToken{}).Where("refresh_token = ?", req.RefreshToken).Update("revoked", true)
	c.JSON(200, model.BaseResponse{Success: true})
}

func generateRefreshToken() string {
	// Implementation of generateRefreshToken function
	// This is a placeholder and should be replaced with the actual implementation
	return ""
}
