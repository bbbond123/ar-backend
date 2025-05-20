package main

import (
	"ar-backend/internal/auth"
	"ar-backend/internal/model"
	server "ar-backend/internal/service"
	"ar-backend/pkg/database"
	"net/http"

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
	// r := router.InitRouter()
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth.NewAuth()
	server := server.NewServer(db)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic("cannot start server")
	}

	// 启动 HTTP 服务
	// r.Run(":8080") // 默认端口 8080

}
