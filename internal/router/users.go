package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// UserRouter 用户路由模块
type UserRouter struct{}

func (UserRouter) Register(r *gin.RouterGroup) {
	user := r.Group("/users")
	{
		user.POST("", controller.CreateUser)
		user.PUT("", controller.UpdateUser)
		user.DELETE(":user_id", controller.DeleteUser)
		user.GET(":user_id", controller.GetUser)
		user.POST("/list", controller.ListUsers)
	}
}

func init() {
	Register(UserRouter{})
}
