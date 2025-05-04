package main

import (
	"ar-backend/internal/router"
	"ar-backend/pkg/database"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	r := gin.Default()

	router.SetupRouter(r)

	r.Run(":8080")
}
