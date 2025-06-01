package auth

import (
	"fmt"
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
			callbackURL = "http://localhost:3000/auth/google/callback"
		}
	}

	// 从环境变量判断是否为生产环境
	isProd := os.Getenv("ENVIRONMENT") == "production"

	fmt.Printf("googleClientId: %s\n", googleClientId)
	fmt.Printf("callbackURL: %s\n", callbackURL)
	fmt.Printf("isProd: %v\n", isProd)
	fmt.Printf("Cookie Secure: %v\n", isProd)
	fmt.Printf("Cookie SameSite: %v\n", func() string {
		if isProd {
			return "None"
		}
		return "Lax"
	}())

	store := sessions.NewCookieStore([]byte(getSessionSecret()))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	// 根据环境设置 SameSite
	if isProd {
		store.Options.SameSite = http.SameSiteNoneMode // 生产环境跨域需要 None
	} else {
		store.Options.SameSite = http.SameSiteLaxMode // 本地开发用 Lax
	}

	// 从环境变量获取域名配置
	if isProd {
		cookieDomain := os.Getenv("COOKIE_DOMAIN")
		if cookieDomain == "" {
			cookieDomain = ".ifoodme.com" // 保持向后兼容
		}
		store.Options.Domain = cookieDomain
	}

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, callbackURL),
	)
}
