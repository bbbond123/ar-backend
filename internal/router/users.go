package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterUsersRoutes 注册设施相关路由
func RegisterUsersRoutes(r *gin.RouterGroup) {
	Users := r.Group("/users")
	{
		Users.POST("", controller.CreateUser)      // 新建
		Users.PUT(":id", controller.UpdateUser)    // 更新
		Users.DELETE(":id", controller.DeleteUser) // 删除
		Users.GET(":id", controller.GetUser)       // 获取单个
		Users.POST("/list", controller.ListUsers)  // 获取列表
	}
}
