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
		// 文件上传（multipart/form-data）
		file.POST("/upload", controller.UploadFile)
		
		// 文件下载
		file.GET("/:file_id/download", controller.DownloadFile)
		
		// S3 连接测试
		file.GET("/test-s3", controller.TestS3Connection)
		
		// 现有的 API
		file.POST("", controller.CreateFile)
		file.PUT("", controller.UpdateFile)
		file.DELETE("/:file_id", controller.DeleteFile)
		file.GET("/:file_id", controller.GetFile)
		file.POST("/list", controller.ListFiles)
	}
}

func init() {
	Register(FileRouter{})
}
