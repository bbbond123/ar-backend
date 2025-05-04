package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// FileRouter 文章路由模块
type FileRouter struct{}

// Register 注册文章路由
func (CommentRouter) FileRegister(r *gin.RouterGroup) {
	files := r.Group("/files")
	{
		files.POST("", controller.CreateFile)
		files.PUT("", controller.UpdateFile)
		files.DELETE(":files_id", controller.DeleteFile)
		files.GET(":files_id", controller.GetFile)
		files.POST("/list", controller.ListFiles)
	}
}

func init() {
	Register(FileRouter{})
}
