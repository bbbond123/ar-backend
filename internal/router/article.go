package router

import (
	"ar-backend/internal/controller"
	"ar-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// ArticleRouter 文章路由模块
type ArticleRouter struct{}

// Register 注册文章路由
func (ArticleRouter) Register(r *gin.RouterGroup) {
	article := r.Group("/articles")
	article.GET(":article_id", controller.GetArticle)
	article.POST("/list", controller.ListArticles)

	aritcleAuth := article.Group("/articles")
	aritcleAuth.Use(middleware.JWTAuth())
	{
		article.POST("", controller.CreateArticle)
		article.PUT("", controller.UpdateArticle)
		article.DELETE(":article_id", controller.DeleteArticle)
	}
}

func init() {
	Register(ArticleRouter{})
}
