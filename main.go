package main

import (
	"ar-backend/internal/model"
	"ar-backend/internal/router"
	"ar-backend/pkg/database"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ar-backend/docs" // Swagger 文档自动生成 (swag init 生成的)
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 初始化数据库
	database.ConnectDatabase()

	// 自动迁移（AutoMigrate会自动创建不存在的表）
	db := database.GetDB()
	db.AutoMigrate(
		&model.Facility{},
		&model.File{},
		&model.Notice{},
		&model.VisitHistory{},
		&model.Language{},
		&model.User{},
		&model.RefreshToken{},
		&model.Store{},
		&model.Menu{},
		&model.Article{},
		&model.Comment{},
		&model.Tag{},
		&model.Tagging{},
	)

	// 初始化 Gin 路由
	r := router.InitRouter()

	// 注册 Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动 HTTP 服务
	r.Run(":8080") // 默认端口 8080
}
