package main

import (
	"ar-backend/internal/model"
	"ar-backend/internal/router"
	"ar-backend/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库

	database.ConnectDatabase()

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

	r := gin.Default()
	router.SetupRouter(r)

	r.Run(":8080")
}
