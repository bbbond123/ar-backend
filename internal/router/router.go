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

	// 公开的认证路由 (不需要JWT验证)
	authPublic := api.Group("/auth")
	{
		authPublic.POST("/login", controller.Login)
		authPublic.POST("/register", controller.Register)
		authPublic.POST("/refresh", controller.RefreshToken)
		authPublic.POST("/logout", controller.RevokeRefreshToken)
		authPublic.POST("/google", controller.GoogleAuth)
	}

	// 需要认证的用户路由
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.JWTAuth())
	{
		authProtected.GET("/user/profile", controller.UserProfile)
	}

	// 注册所有模块路由
	for _, rr := range routeRegisters {
		rr.Register(api)
	}

	return r
}
