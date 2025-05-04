package router

import (
	"ar-backend/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/ping", controller.Ping)
	r.GET("/test/db", controller.TestDB)
}
