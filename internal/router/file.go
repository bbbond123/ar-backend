package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// FileRouter 文件路由模块
type FileRouter struct{}

// Register 注册文件路由
func (FileRouter) Register(r *gin.RouterGroup) {
	file := r.Group("/files")
	{
		file.POST("", controller.CreateFile)
		file.PUT("", controller.UpdateFile)
		file.DELETE(":file_id", controller.DeleteFile)
		file.GET(":file_id", controller.GetFile)
		file.POST("/list", controller.ListFiles)
	}
}

func init() {
	Register(FileRouter{})
}
