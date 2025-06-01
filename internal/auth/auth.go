package auth

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	MaxAge = 86400 * 30
)

// getSessionSecret 从环境变量获取 Session 密钥
func getSessionSecret() string {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		// 如果没有设置专门的session secret，使用JWT_SECRET
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Fatal("SESSION_SECRET or JWT_SECRET environment variable is required")
		}
		return jwtSecret + "_session"
	}
	return secret
}

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	// 从环境变量获取回调 URL
	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")
	if callbackURL == "" {
		// 根据环境设置默认回调URL
		if os.Getenv("ENVIRONMENT") == "production" {
			callbackURL = "https://www.ifoodme.com/api/auth/google/callback"
		} else {
			callbackURL = "http://localhost:3000/api/auth/google/callback"
		}
	}

	// 从环境变量判断是否为生产环境
	isProd := os.Getenv("ENVIRONMENT") == "production"

	log.Printf("=== OAuth Configuration ===")
	log.Printf("Environment: %s", os.Getenv("ENVIRONMENT"))
	log.Printf("Google Client ID: %s", googleClientId)
	log.Printf("Google Callback URL: %s", callbackURL)
	log.Printf("Is Production: %v", isProd)
	log.Printf("Cookie Secure: %v", isProd)
	log.Printf("Cookie SameSite: %v", func() string {
		if isProd {
			return "None"
		}
		return "Lax"
	}())
	log.Printf("========================")

	store := sessions.NewCookieStore([]byte(getSessionSecret()))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	// 根据环境设置 SameSite
	if isProd {
		store.Options.SameSite = http.SameSiteNoneMode // 生产环境跨域需要 None
		// 生产环境设置域名
		cookieDomain := os.Getenv("COOKIE_DOMAIN")
		if cookieDomain == "" {
			cookieDomain = ".ifoodme.com" // 保持向后兼容
		}
		store.Options.Domain = cookieDomain
		log.Printf("Production Cookie Domain: %s", cookieDomain)
	} else {
		store.Options.SameSite = http.SameSiteLaxMode // 本地开发用 Lax
		log.Printf("Development mode - no domain restriction")
	}

	gothic.Store = store

	// 设置 Google OAuth 配置
	googleProvider := google.New(googleClientId, googleClientSecret, callbackURL)
	googleProvider.SetPrompt("select_account") // 强制显示 Google 账号选择界面
	goth.UseProviders(googleProvider)
}
