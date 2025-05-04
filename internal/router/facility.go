package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterFacilityRoutes 注册设施相关路由
func RegisterFacilityRoutes(r *gin.RouterGroup) {
	facility := r.Group("/facilities")
	{
		facility.POST("", controller.CreateFacility)      // 新建
		facility.PUT(":id", controller.UpdateFacility)    // 更新
		facility.DELETE(":id", controller.DeleteFacility) // 删除
		facility.GET(":id", controller.GetFacility)       // 获取单个
		facility.POST("/list", controller.ListFacilities) // 获取列表
	}
}
