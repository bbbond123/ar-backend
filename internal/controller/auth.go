package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"ar-backend/internal/model"
	"ar-backend/pkg/database"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
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
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: err.Error(), Code: 400})
		return
	}
	db := database.GetDB()
	var user model.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(401, model.BaseResponse{Success: false, ErrMessage: "用户不存在", Code: 401})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(401, model.BaseResponse{Success: false, ErrMessage: "密码错误", Code: 401})
		return
	}
	accessToken, err := generateAccessToken(user.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "Token生成失败", Code: 500})
		return
	}
	refreshToken, err := generateRefreshTokenJWT(user.UserID)
	if err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "RefreshToken生成失败", Code: 500})
		return
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := db.Create(&model.RefreshToken{
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Revoked:      false,
	}).Error; err != nil {
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "RefreshToken存储失败", Code: 500})
		return
	}
	c.JSON(200, model.Response[model.AuthResponse]{
		Success: true,
		Code:    200,
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
	fmt.Printf("=== Google Token验证开始 ===\n")
	fmt.Printf("接收到的ID Token长度: %d\n", len(idToken))
	fmt.Printf("ID Token前50字符: %s...\n", func() string {
		if len(idToken) > 50 {
			return idToken[:50]
		}
		return idToken
	}())

	// 支持多平台的Client ID验证
	validClientIDs := []string{
		"680314480886-ugffmjjjdfdg1a98g5ami0sa9f10pbbn.apps.googleusercontent.com", // Android
		"680314480886-o8el90n41jc8g14qvu526a6iuflucpiu.apps.googleusercontent.com", // iOS
		"680314480886-87foecji3cgqu9vqt85eh7ua6r6bnn9s.apps.googleusercontent.com", // Web
		os.Getenv("GOOGLE_CLIENT_ID"), // 环境变量中的Client ID
	}

	fmt.Printf("允许的Client IDs: %v\n", validClientIDs)

	tokenInfoURL := "https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken
	fmt.Printf("请求Google Token验证URL: %s\n", tokenInfoURL[:100]+"...")

	resp, err := http.Get(tokenInfoURL)
	if err != nil {
		fmt.Printf("Google Token验证请求失败: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("Google API响应状态码: %d\n", resp.StatusCode)

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Google API错误响应: %s\n", string(bodyBytes))
		return nil, fmt.Errorf("invalid google id_token, status: %d", resp.StatusCode)
	}

	var info GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		fmt.Printf("解析Google用户信息失败: %v\n", err)
		return nil, err
	}

	fmt.Printf("获取到的用户信息 - Email: %s, Name: %s, Sub: %s\n", info.Email, info.Name, info.Sub)

	// 验证token的audience (client_id)是否在我们的允许列表中
	tokenResp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
	if err != nil {
		fmt.Printf("二次Token验证请求失败: %v\n", err)
		return nil, err
	}
	defer tokenResp.Body.Close()

	var tokenInfo struct {
		Aud           string `json:"aud"` // client_id
		Sub           string `json:"sub"` // user_id
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Iss           string `json:"iss"` // issuer
		Exp           string `json:"exp"` // expiration
	}

	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenInfo); err != nil {
		fmt.Printf("解析Token详细信息失败: %v\n", err)
		return nil, err
	}

	fmt.Printf("Token详细信息 - Aud(Client ID): %s, Iss: %s, Exp: %s\n", tokenInfo.Aud, tokenInfo.Iss, tokenInfo.Exp)
	fmt.Printf("Email验证状态: %s\n", tokenInfo.EmailVerified)

	// 检查client_id是否在允许列表中
	validClientID := false
	for _, clientID := range validClientIDs {
		if clientID != "" && tokenInfo.Aud == clientID {
			validClientID = true
			fmt.Printf("✅ Google token验证成功，匹配的client_id: %s\n", clientID)
			break
		}
	}

	if !validClientID {
		fmt.Printf("❌ 无效的client_id: %s，不在允许列表中\n", tokenInfo.Aud)
		fmt.Printf("允许的Client IDs: %v\n", validClientIDs)
		return nil, fmt.Errorf("invalid client_id: %s", tokenInfo.Aud)
	}

	// 返回验证后的用户信息
	result := &GoogleUserInfo{
		Sub:           tokenInfo.Sub,
		Email:         tokenInfo.Email,
		EmailVerified: tokenInfo.EmailVerified,
		Name:          tokenInfo.Name,
		Picture:       tokenInfo.Picture,
	}

	fmt.Printf("=== Google Token验证成功 ===\n")
	fmt.Printf("返回用户信息: %+v\n", result)

	return result, nil
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
	fmt.Printf("\n=== TTTTT GoogleAuth POST 开始 ===\n")
	fmt.Printf("请求IP: %s\n", c.ClientIP())
	fmt.Printf("User-Agent: %s\n", c.GetHeader("User-Agent"))
	fmt.Printf("X-App-Platform: %s\n", c.GetHeader("X-App-Platform"))

	var req model.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("❌ 参数绑定失败: %v\n", err)
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "参数错误: " + err.Error()})
		return
	}
	fmt.Printf("111222333 Token长度: %d\n", req)
	fmt.Printf("接收到ID Token长度: %d\n", len(req.IdToken))

	// 验证id_token，获取Google用户信息
	userInfo, err := getGoogleUserInfo(req.IdToken)
	if err != nil {
		fmt.Printf("❌ Google token验证失败: %v\n", err)
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "Google token无效"})
		return
	}
	if userInfo.Email == "" || userInfo.Sub == "" {
		fmt.Printf("❌ Google用户信息不完整 - Email: %s, Sub: %s\n", userInfo.Email, userInfo.Sub)
		c.JSON(400, model.BaseResponse{Success: false, ErrMessage: "Google用户信息不完整"})
		return
	}

	fmt.Printf("✅ Google用户信息验证成功 - Email: %s, Sub: %s\n", userInfo.Email, userInfo.Sub)

	db := database.GetDB()
	var user model.User
	if err := db.Where("google_id = ?", userInfo.Sub).First(&user).Error; err == nil {
		fmt.Printf("✅ 找到已存在用户 - UserID: %d, Email: %s\n", user.UserID, user.Email)
		// 已存在，直接登录
		if user.Status != "active" {
			fmt.Printf("更新用户状态为active - UserID: %d\n", user.UserID)
			user.Status = "active"
			db.Save(&user)
		}
	} else {
		fmt.Printf("🆕 创建新用户 - Email: %s, GoogleID: %s\n", userInfo.Email, userInfo.Sub)
		// 不存在，注册
		user = model.User{
			Email:    userInfo.Email,
			GoogleID: userInfo.Sub,
			Name:     userInfo.Name,
			Provider: "google",
			Status:   "active",
		}
		if err := db.Create(&user).Error; err != nil {
			fmt.Printf("❌ 用户注册失败: %v\n", err)
			c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "用户注册失败: " + err.Error()})
			return
		}
		fmt.Printf("✅ 新用户创建成功 - UserID: %d\n", user.UserID)
	}

	fmt.Printf("开始生成Token - UserID: %d\n", user.UserID)

	accessToken, err := generateAccessToken(user.UserID)
	if err != nil {
		fmt.Printf("❌ AccessToken生成失败: %v\n", err)
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "Token生成失败"})
		return
	}

	fmt.Printf("✅ AccessToken生成成功，长度: %d\n", len(accessToken))

	refreshToken, err := generateRefreshTokenJWT(user.UserID)
	if err != nil {
		fmt.Printf("❌ RefreshToken生成失败: %v\n", err)
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "RefreshToken生成失败"})
		return
	}

	fmt.Printf("✅ RefreshToken生成成功，长度: %d\n", len(refreshToken))

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := db.Create(&model.RefreshToken{
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		Revoked:      false,
	}).Error; err != nil {
		fmt.Printf("❌ RefreshToken存储失败: %v\n", err)
		c.JSON(500, model.BaseResponse{Success: false, ErrMessage: "RefreshToken存储失败"})
		return
	}

	fmt.Printf("✅ RefreshToken存储成功，过期时间: %v\n", expiresAt)

	fmt.Printf("=== GoogleAuth POST 成功完成 ===\n")
	fmt.Printf("用户: %s (ID: %d) 登录成功\n\n", user.Email, user.UserID)

	c.JSON(200, model.Response[model.AuthResponse]{
		Success: true,
		Data: model.AuthResponse{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// BeginGoogleAuth godoc
// @Summary 开始 Google OAuth 认证
// @Description 重定向到 Google 登录页面
// @Tags Auth
// @Accept json
// @Produce json
// @Param redirect query string false "重定向URL，支持深度链接(如: travelview://google-auth-callback)"
// @Success 302 {string} string "重定向到 Google 登录页面"
// @Router /api/auth/google [get]
func BeginGoogleAuth(c *gin.Context) {
	fmt.Printf("\n=== BeginGoogleAuth GET 开始 ===\n")
	fmt.Printf("请求IP: %s\n", c.ClientIP())
	fmt.Printf("User-Agent: %s\n", c.GetHeader("User-Agent"))
	fmt.Printf("Referer: %s\n", c.GetHeader("Referer"))

	provider := "google"
	ctx := context.WithValue(context.Background(), "provider", provider)
	r := c.Request.WithContext(ctx)
	w := c.Writer

	// 获取前端传递的 redirect 参数，支持深度链接
	// 支持两种参数名：redirect 和 redirect_uri
	redirectURL := c.Query("redirect")
	if redirectURL == "" {
		redirectURL = c.Query("redirect_uri")
	}
	fmt.Printf("接收到的redirect参数: %s\n", redirectURL)
	fmt.Printf("1111原始查询字符串: %s\n", c.Request.URL.RawQuery)

	if redirectURL != "" {
		fmt.Printf("BeginGoogleAuth - 处理redirect参数: %s\n", redirectURL)
		// 将 redirect URL 保存到 session 中，以便在回调时使用
		session, err := gothic.Store.Get(r, "oauth_session")
		if err == nil {
			session.Values["redirect_url"] = redirectURL
			err = session.Save(r, w)
			if err != nil {
				fmt.Printf("❌ BeginGoogleAuth - 保存session失败: %v\n", err)
			} else {
				fmt.Printf("✅ BeginGoogleAuth - 成功保存redirect URL到session: %s\n", redirectURL)
			}
		} else {
			fmt.Printf("❌ BeginGoogleAuth - 获取session失败: %v\n", err)
		}
	} else {
		fmt.Printf("无redirect参数，将使用默认前端URL\n")
	}

	fmt.Printf("准备重定向到Google OAuth页面...\n")
	fmt.Printf("=== BeginGoogleAuth GET 结束 ===\n\n")

	gothic.BeginAuthHandler(w, r)
}

// GoogleAuthCallback godoc
// @Summary Google OAuth 回调处理
// @Description 处理 Google 登录回调，创建或更新用户，并返回 JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 302 {string} string "重定向到前端页面，并携带 token"
// @Router /api/auth/google/callback [get]
func GoogleAuthCallback(c *gin.Context) {
	fmt.Printf("\n=== GoogleAuthCallback 开始 ===\n")
	fmt.Printf("请求IP: %s\n", c.ClientIP())
	fmt.Printf("User-Agent: %s\n", c.GetHeader("User-Agent"))
	fmt.Printf("查询参数: %s\n", c.Request.URL.RawQuery)

	provider := "google"
	ctx := context.WithValue(context.Background(), "provider", provider)
	r := c.Request.WithContext(ctx)
	w := c.Writer

	fmt.Printf("开始完成Google OAuth认证...\n")
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Printf("❌ Google OAuth认证失败: %v\n", err)
		c.String(http.StatusUnauthorized, "auth error: %v", err)
		return
	}

	fmt.Printf("✅ Google OAuth认证成功\n")
	fmt.Printf("Google用户信息 - Email: %s, Name: %s, UserID: %s\n", user.Email, user.Name, user.UserID)
	fmt.Printf("Avatar: %s, Provider: %s\n", user.AvatarURL, user.Provider)

	db := database.GetDB()
	var userInDB model.User
	err = db.Where("email = ?", user.Email).First(&userInDB).Error
	if err == gorm.ErrRecordNotFound {
		fmt.Printf("🆕 创建新用户 - Email: %s\n", user.Email)
		userInDB = model.User{
			Email:    user.Email,
			GoogleID: user.UserID,
			Name:     user.Name,
			Avatar:   user.AvatarURL,
			Provider: "google",
			Status:   "active",
		}
		err = db.Create(&userInDB).Error
		if err != nil {
			fmt.Printf("❌ 创建用户失败: %v\n", err)
			c.String(http.StatusInternalServerError, "Could not create user")
			return
		}
		fmt.Printf("✅ 新用户创建成功 - UserID: %d\n", userInDB.UserID)
	} else if err == nil {
		fmt.Printf("✅ 找到已存在用户 - UserID: %d\n", userInDB.UserID)
		err = db.Model(&userInDB).Updates(map[string]interface{}{
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

	// 从环境变量获取域名配置
	isProd := os.Getenv("ENVIRONMENT") == "production"
	if isProd {
		cookieDomain := os.Getenv("COOKIE_DOMAIN")
		if cookieDomain == "" {
			cookieDomain = ".ifoodme.com" // 保持向后兼容
		}
		cookie.Domain = cookieDomain
		fmt.Printf("设置生产环境Cookie - Domain: %s, Secure: true\n", cookieDomain)
	} else {
		cookie.Secure = false
		fmt.Printf("设置开发环境Cookie - Secure: false\n")
	}

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
			fmt.Printf("✅ 从session获取到redirect URL: %s\n", frontendURL)
		} else {
			fmt.Printf("❌ session中没有找到redirect_url或值为空\n")
			fmt.Printf("redirect_url值类型: %T, 值: %v\n", session.Values["redirect_url"], session.Values["redirect_url"])
		}
	} else {
		fmt.Printf("❌ 获取session失败: %v\n", err)
	}

	// 2. 如果没有 redirect 参数，使用环境变量
	if frontendURL == "" {
		frontendURL = os.Getenv("FRONTEND_URL")
		fmt.Printf("使用环境变量FRONTEND_URL: %s\n", frontendURL)
	}

	// 3. 如果环境变量也没有，使用默认地址
	if frontendURL == "" {
		frontendURL = "https://www.ifoodme.com"
		fmt.Printf("使用默认前端URL: %s\n", frontendURL)
	}

	// 检查是否是React Native应用的深度链接
	if strings.HasPrefix(frontendURL, "travelview://") {
		fmt.Printf("🔗 检测到React Native深度链接\n")
		// React Native 深度链接，构造参数
		deepLink := frontendURL + "?token=" + url.QueryEscape(tokenString) +
			"&user_id=" + fmt.Sprintf("%d", userInDB.UserID) +
			"&email=" + url.QueryEscape(userInDB.Email) +
			"&name=" + url.QueryEscape(userInDB.Name)

		if userInDB.Avatar != "" {
			deepLink += "&avatar=" + url.QueryEscape(userInDB.Avatar)
		}

		fmt.Printf("构造的深度链接: %s\n", deepLink)
		fmt.Printf("=== 重定向到React Native App ===\n\n")

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

	fmt.Printf("最终重定向URL: %s\n", frontendURL)
	fmt.Printf("=== GoogleAuthCallback 成功完成 ===\n")
	fmt.Printf("用户: %s (ID: %d) 通过OAuth登录成功\n\n", userInDB.Email, userInDB.UserID)

	// 重定向到前端首页
	c.Redirect(http.StatusFound, frontendURL)
}
