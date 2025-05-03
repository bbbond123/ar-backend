// @title My API
// @version 1.0
// @description API for ar-backend
// @BasePath /api/v1
package main

import (
	"ar-backend/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.SetupRouter(r)

	r.Run(":8080")
}
