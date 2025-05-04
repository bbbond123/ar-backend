package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// NoticeRouter
type NoticeRouter struct{}

func (CommentRouter) NoticeRegister(r *gin.RouterGroup) {
	notice := r.Group("/notices")
	{
		notice.POST("", controller.CreateNotice)
		notice.PUT("", controller.UpdateNotice)
		notice.DELETE(":notice_id", controller.DeleteNotice)
		notice.GET(":notice_id", controller.GetNotice)
		notice.POST("/list", controller.ListNotices)
	}
}

func init() {
	Register(NoticeRouter{})
}
