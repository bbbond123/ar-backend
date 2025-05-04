package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

// ArticleRouter 文章路由模块
type ArticleRouter struct{}

// Register 注册文章路由
func (ArticleRouter) Register(r *gin.RouterGroup) {
	article := r.Group("/articles")
	{
		article.POST("", controller.CreateArticle)
		article.PUT("", controller.UpdateArticle)
		article.DELETE(":article_id", controller.DeleteArticle)
		article.GET(":article_id", controller.GetArticle)
		article.POST("/list", controller.ListArticles)
	}
}

func init() {
	Register(ArticleRouter{})
}
