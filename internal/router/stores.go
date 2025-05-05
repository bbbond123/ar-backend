package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// StoreRouter Refresh Token 路由模块
type StoreRouter struct{}

// Register 注册 Refresh Token 路由
func (StoreRouter) Register(r *gin.RouterGroup) {
	Store := r.Group("/stores")
	{
		Store.POST("", controller.CreateStore)
		Store.PUT("", controller.UpdateStore)
		Store.DELETE(":token_id", controller.DeleteStore)
		Store.GET(":token_id", controller.GetStore) //
		Store.POST("/list", controller.ListStores)
		// Store.GET(":store_id/tags", controller.GetTagsByStore)
		// Store.POST(":store_id/tags", controller.AddTagToStore) //
		// Store.DELETE(":store_id/tags/:tag_id", controller.RemoveTagFromStore)
	}
}

func init() {
	Register(StoreRouter{})
}
