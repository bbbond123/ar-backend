package main

import (
	"ar-backend/internal/auth"
	"ar-backend/internal/model"
	server "ar-backend/internal/service"
	"ar-backend/pkg/database"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "ar-backend/docs" // Swagger 文档自动生成 (swag init 生成的)

	"github.com/joho/godotenv"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	fmt.Println("🚀 启动 AR Backend 服务...")

	// 首先加载 .env 文件
	fmt.Println("📋 正在加载环境配置...")
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ 未找到 .env 文件，使用系统环境变量")
	} else {
		fmt.Println("✅ .env 文件加载成功")
	}

	// 初始化数据库
	fmt.Println("📊 正在连接数据库...")
	database.ConnectDatabase()
	fmt.Println("✅ 数据库连接成功")

	// 自动迁移（AutoMigrate会自动创建不存在的表）
	fmt.Println("🔄 正在进行数据库迁移...")
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
	fmt.Println("✅ 数据库迁移完成")

	// 初始化示例用户数据
	fmt.Println("👥 正在初始化用户数据...")
	server.InitializeSampleUsers()

	// 初始化认证
	fmt.Println("🔐 正在初始化认证模块...")
	auth.NewAuth()
	fmt.Println("✅ 认证模块初始化完成")

	// 启动服务器
	fmt.Println("🌐 正在启动HTTP服务器...")
	serverInstance := server.NewServer(db)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("✅ 服务器启动成功!\n")
	fmt.Printf("🌐 服务地址: http://localhost:%s\n", port)
	fmt.Printf("📖 API文档: http://localhost:%s/swagger/index.html\n", port)
	fmt.Printf("💚 健康检查: http://localhost:%s/api/health\n", port)

	err = serverInstance.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("❌ 服务器启动失败: %v\n", err)
		panic("cannot start server")
	}
}
