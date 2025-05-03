package controller

import "github.com/gin-gonic/gin"

// Ping godoc
// @Summary Ping API
// @Produce json
// @Success 200 {string} string "pong"
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

// swagger docs
// 安装 swaggo 生成
// swag init -g cmd/main.go
