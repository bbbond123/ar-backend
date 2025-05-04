package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// MenusRouter 文章路由模块
type MenusRouter struct{}

// Register 注册文章路由
func (CommentRouter) MenusRegister(r *gin.RouterGroup) {
	Menus := r.Group("/menus")
	{
		Menus.POST("", controller.CreateMenu)
		Menus.PUT("", controller.UpdateMenu)
		Menus.DELETE(":menus_id", controller.DeleteMenu)
		Menus.GET(":menus_id", controller.GetMenu)
		Menus.POST("/list", controller.ListMenus)
	}
}

func init() {
	Register(CommentRouter{})
}
