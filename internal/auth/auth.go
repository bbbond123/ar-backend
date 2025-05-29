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
	key    = "randomString"
	MaxAge = 86400 * 30
)

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	// 从环境变量获取回调 URL，如果没有则使用默认值
	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")
	if callbackURL == "" {
		callbackURL = "https://api.ifoodme.com/api/auth/google/callback"
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

	store := sessions.NewCookieStore([]byte(key))
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

	// 生产环境下不设置域名，让浏览器自动处理
	// if isProd {
	// 	store.Options.Domain = ".ifoodme.com" // 允许子域名共享
	// }

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, callbackURL),
	)
}
