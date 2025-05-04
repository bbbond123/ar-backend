package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// CommentRouter 评论路由模块
type CommentRouter struct{}

// Register 注册评论路由
func (CommentRouter) Register(r *gin.RouterGroup) {
	comment := r.Group("/comments")
	{
		comment.POST("", controller.CreateComment)
		comment.PUT("", controller.UpdateComment)
		comment.DELETE(":comment_id", controller.DeleteComment)
		comment.GET(":comment_id", controller.GetComment)
		comment.POST("/list", controller.ListComments)
	}
}

func init() {
	Register(CommentRouter{})
}
