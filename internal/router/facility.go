package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// FacilityRouter 设施路由模块
type FacilityRouter struct{}

// Register 注册设施路由
func (FacilityRouter) Register(r *gin.RouterGroup) {
	facility := r.Group("/facilities")
	{
		facility.POST("", controller.CreateFacility)
		facility.PUT(":id", controller.UpdateFacility)
		facility.DELETE(":id", controller.DeleteFacility)
		facility.GET(":id", controller.GetFacility)
		facility.POST("/list", controller.ListFacilities)
	}
}

func init() {
	Register(FacilityRouter{})
}
