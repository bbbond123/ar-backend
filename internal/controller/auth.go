package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"ar-backend/internal/model"
	"ar-backend/pkg/database"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// getJWTSecret 从环境变量获取 JWT 密钥
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	return []byte(secret)
}

// getRefreshSecret 从环境变量获取 Refresh Token 密钥
func getRefreshSecret() []byte {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		// 如果没有设置专门的refresh secret，使用JWT_SECRET + 后缀
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Fatal("JWT_SECRET environment variable is required")
		}
		return []byte(jwtSecret + "_refresh")
	}
	return []byte(secret)
}

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
	accessToken, err := generateAccessToken(user.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "Token生成失败"})
		return
	}
	refreshToken, err := generateRefreshTokenJWT(user.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "RefreshToken生成失败"})
		return
	}
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
	var user model.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err == nil {
		// 已存在该邮箱
		if user.Status == "pending" {
			// 未激活，更新验证码和过期时间
			verifyCode := fmt.Sprintf("%04d", rand.Intn(10000))
			verifyExpire := time.Now().Add(10 * time.Minute)
			user.VerifyCode = verifyCode
			user.VerifyCodeExpire = &verifyExpire
			if err := db.Save(&user).Error; err != nil {
				c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "验证码更新失败: " + err.Error()})
				return
			}
			// TODO: 发送验证码到邮箱 user.Email，内容为 verifyCode
			// sendVerifyCodeToEmail(user.Email, verifyCode)
			c.JSON(200, model.BaseResponse{Success: true, ErrMessage: "验证码已重新发送，请查收邮箱"})
			return
		} else {
			c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "邮箱已被注册"})
			return
		}
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "密码加密失败"})
		return
	}

	// 生成4位验证码
	verifyCode := fmt.Sprintf("%04d", rand.Intn(10000))
	verifyExpire := time.Now().Add(10 * time.Minute)

	user = model.User{
		Email:            req.Email,
		Password:         string(hashedPwd),
		Provider:         "email",
		Status:           "pending", // 注册后状态为pending，待激活
		VerifyCode:       verifyCode,
		VerifyCodeExpire: &verifyExpire,
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	// TODO: 发送验证码到邮箱 user.Email，内容为 verifyCode
	// sendVerifyCodeToEmail(user.Email, verifyCode)

	accessToken, err := generateAccessToken(user.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "Token生成失败"})
		return
	}

	refreshToken, err := generateRefreshTokenJWT(user.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "RefreshToken生成失败"})
		return
	}

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

// 生成短时access token（15分钟）
func generateAccessToken(userID int) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
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
	return token.SignedString(getJWTSecret())
}

// 生成长时refresh token（7天）
func generateRefreshTokenJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	type RefreshClaims struct {
		UserID int `json:"user_id"`
		jwt.RegisteredClaims
	}
	claims := &RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getRefreshSecret())
}

// RefreshToken godoc
// @Summary 刷新Access Token
// @Description 使用Refresh Token刷新Access Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body model.RefreshTokenRequest true "刷新Token请求"
// @Success 200 {object} model.Response[model.RefreshTokenResponse]
// @Failure 400 {object} model.BaseResponse
// @Failure 401 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	// 1. 校验refreshToken格式和签名
	refreshTokenStr := req.RefreshToken
	type RefreshClaims struct {
		UserID int `json:"user_id"`
		jwt.RegisteredClaims
	}
	claims := &RefreshClaims{}
	token, err := jwt.ParseWithClaims(refreshTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return getRefreshSecret(), nil
	})
	if err != nil || !token.Valid {
		c.JSON(401, model.BaseResponse{Success: false, ErrMessage: "refresh token无效"})
		return
	}
	// 2. 查库校验refreshToken是否存在且未撤销且未过期
	db := database.GetDB()
	var dbToken model.RefreshToken
	if err := db.Where("refresh_token = ? AND revoked = false AND expires_at > ?", refreshTokenStr, time.Now()).First(&dbToken).Error; err != nil {
		c.JSON(401, model.BaseResponse{Success: false, ErrMessage: "refresh token无效或已过期"})
		return
	}
	// 3. 生成新的access token
	accessToken, err := generateAccessToken(claims.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "Token生成失败"})
		return
	}
	c.JSON(200, model.Response[model.RefreshTokenResponse]{Success: true, Data: model.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
	}})
}

// RevokeRefreshToken godoc
// @Summary 登出
// @Description 使refresh token失效（登出）
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body model.RevokeTokenRequest true "登出请求"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/auth/logout [post]
func RevokeRefreshToken(c *gin.Context) {

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

	// var req model.RevokeTokenRequest
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(400, model.BaseResponse{Success: false, ErrMessage: err.Error()})
	// 	return
	// }
	// db := database.GetDB()
	// result := db.Model(&model.RefreshToken{}).Where("refresh_token = ? AND revoked = false", req.RefreshToken).Update("revoked", true)
	// if result.Error != nil {
	// 	c.JSON(500, model.BaseResponse{Success: false, ErrMessage: result.Error.Error()})
	// 	return
	// }
	// if result.RowsAffected == 0 {
	// 	c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "无效或已撤销的refresh token"})
	// 	return
	// }
	// c.JSON(200, model.BaseResponse{Success: true})
}

// GoogleUserInfo 用于解析Google返回的用户信息
type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func getGoogleUserInfo(idToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid google id_token")
	}
	var info GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GoogleAuth godoc
// @Summary Google社交登录/注册
// @Description Google社交登录/注册
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body model.GoogleAuthRequest true "Google登录请求"
// @Success 200 {object} model.Response[model.AuthResponse]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/auth/google [post]
func GoogleAuth(c *gin.Context) {
	var req model.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "参数错误: " + err.Error()})
		return
	}

	// 验证id_token，获取Google用户信息
	userInfo, err := getGoogleUserInfo(req.IdToken)
	if err != nil {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "Google token无效"})
		return
	}
	if userInfo.Email == "" || userInfo.Sub == "" {
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "Google用户信息不完整"})
		return
	}

	db := database.GetDB()
	var user model.User
	if err := db.Where("google_id = ?", userInfo.Sub).First(&user).Error; err == nil {
		// 已存在，直接登录
		if user.Status != "active" {
			user.Status = "active"
			db.Save(&user)
		}
	} else {
		// 不存在，注册
		user = model.User{
			Email:    userInfo.Email,
			GoogleID: userInfo.Sub,
			Name:     userInfo.Name,
			Provider: "google",
			Status:   "active",
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "用户注册失败: " + err.Error()})
			return
		}
	}

	accessToken, err := generateAccessToken(user.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "Token生成失败"})
		return
	}

	refreshToken, err := generateRefreshTokenJWT(user.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "RefreshToken生成失败"})
		return
	}

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
