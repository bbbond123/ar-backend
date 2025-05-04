package router

import (
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// API分组
	api := r.Group("/api")
	{
		RegisterFacilityRoutes(api) // 注册设施路由
	}

	return r
}
