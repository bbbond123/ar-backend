package main

import (
	"ar-backend/internal/router"
	"ar-backend/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库

	database.ConnectDatabase()
	r := gin.Default()
	router.SetupRouter(r)

	r.Run(":8080")
}
