package controller

import (
	"ar-backend/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping godoc
// @Summary Ping API
// @Produce json
// @Success 200 {string} string "pong"
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func TestDB(c *gin.Context) {
	var result map[string]interface{}

	// 这里用 facilities 表测试（你可以换成其他表）
	if err := database.GetDB().Table("facilities").First(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": result})
}

// swagger docs
// 安装 swaggo 生成
// swag init -g cmd/main.go
