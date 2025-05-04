package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// TagRouter 标签路由模块
type TagRouter struct{}

// Register 注册标签相关路由
func (TagRouter) Register(r *gin.RouterGroup) {
	tags := r.Group("/tags")
	{
		tags.POST("", controller.CreateTag)          // 新建标签
		tags.PUT("", controller.UpdateTag)           // 更新标签
		tags.DELETE(":tag_id", controller.DeleteTag) // 删除标签
		tags.GET(":tag_id", controller.GetTag)       // 获取单个标签
		tags.POST("/list", controller.ListTags)      // 标签分页列表
	}
}

func init() {
	Register(TagRouter{})
}
