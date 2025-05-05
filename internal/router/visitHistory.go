package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// VisitHistoryRouter 访问记录路由模块
type VisitHistoryRouter struct{}

// Register 注册访问记录路由
func (VisitHistoryRouter) Register(r *gin.RouterGroup) {
	visitHistory := r.Group("/visit_history")
	{
		visitHistory.POST("", controller.CreateVisitHistory)
		visitHistory.PUT("", controller.UpdateVisitHistory)
		visitHistory.DELETE(":history_id", controller.DeleteVisitHistory)
		visitHistory.GET(":history_id", controller.GetVisitHistory)
		visitHistory.POST("/list", controller.ListVisitHistories)
	}
}

func init() {
	Register(VisitHistoryRouter{})
}
