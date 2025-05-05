package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// LanguageRouter 语言路由模块
type LanguageRouter struct{}

// Register 注册语言相关路由
func (LanguageRouter) Register(r *gin.RouterGroup) {
	languages := r.Group("/languages")
	{
		languages.POST("", controller.CreateLanguage)               // 新建语言
		languages.PUT("", controller.UpdateLanguage)                // 更新语言
		languages.DELETE(":language_id", controller.DeleteLanguage) // 删除语言
		languages.GET(":language_id", controller.GetLanguage)       // 获取单个语言
		languages.POST("/list", controller.ListLanguages)           // 获取语言分页列表
	}
}

func init() {
	Register(LanguageRouter{})
}
