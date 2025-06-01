package router

import (
	"ar-backend/internal/controller"
	"ar-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// ArticleRouter 文章路由模块
type ArticleRouter struct{}

// Register 注册文章路由
func (ArticleRouter) Register(api *gin.RouterGroup) {
	article := api.Group("/articles")
	
	// 公开访问的路由
	api.GET("/articles/:article_id", controller.GetArticle)
	api.POST("/articles/list", controller.ListArticles)
	
	// 带图片上传的文章创建（暂时不需要认证，方便测试）
	article.POST("/with-image", controller.CreateArticleWithImage)

	// 需要认证的路由
	aritcleAuth := article.Group("/articles")
	aritcleAuth.Use(middleware.JWTAuth())
	{
		article.POST("", controller.CreateArticle)
		article.PUT("", controller.UpdateArticle)
		article.DELETE("/:article_id", controller.DeleteArticle)
	}
}

func init() {
	Register(ArticleRouter{})
}
