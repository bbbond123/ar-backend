package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// MenuRouter 菜单路由模块
type MenuRouter struct{}

// Register 注册菜单路由
func (MenuRouter) Register(r *gin.RouterGroup) {
	menu := r.Group("/menus")
	{
		menu.POST("", controller.CreateMenu)
		menu.PUT("", controller.UpdateMenu)
		menu.DELETE(":menu_id", controller.DeleteMenu)
		menu.GET(":menu_id", controller.GetMenu)
		menu.POST("/list", controller.ListMenus)
	}
}

func init() {
	Register(MenuRouter{})
}
