package router

import (
	"ar-backend/internal/controller"
	"ar-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RouteRegister 路由注册器接口
type RouteRegister interface {
	Register(r *gin.RouterGroup)
}

var routeRegisters []RouteRegister

// Register 注册路由模块
func Register(rr RouteRegister) {
	routeRegisters = append(routeRegisters, rr)
}

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	api.POST("/login", controller.Login)
	api.POST("/register", controller.Register)
	api.POST("/refresh", controller.RefreshToken)
	// api.POST("/logout", controller.RevokeRefreshToken)

	// 注册所有模块路由
	for _, rr := range routeRegisters {
		rr.Register(api)
	}

	// 注册auth路由
	auth := api.Group("/auth")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/user/profile", controller.UserProfile)
		// 其他需要登录的接口
	}

	return r
}
