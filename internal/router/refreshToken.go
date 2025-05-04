package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// RefreshTokenRouter Refresh Token 路由模块
type RefreshTokenRouter struct{}

// Register 注册 Refresh Token 路由
func (RefreshTokenRouter) Register(r *gin.RouterGroup) {
	refreshToken := r.Group("/refresh_tokens")
	{
		refreshToken.POST("", controller.CreateRefreshToken)
		refreshToken.PUT("", controller.UpdateRefreshToken)
		refreshToken.DELETE(":token_id", controller.DeleteRefreshToken)
		refreshToken.GET(":token_id", controller.GetRefreshToken)
		refreshToken.POST("/list", controller.ListRefreshTokens)
	}
}

func init() {
	Register(RefreshTokenRouter{})
}
