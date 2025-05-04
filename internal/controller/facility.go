package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateFacility(c *gin.Context) {
	var facility model.Facility
	if err := c.ShouldBindJSON(&facility); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := database.GetDB()
	if err := db.Create(&facility).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, facility)
}

func ListFacilities(c *gin.Context) {
	db := database.GetDB()
	var facilities []model.Facility
	db.Find(&facilities)
	c.JSON(http.StatusOK, facilities)
}

func GetFacility(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var facility model.Facility
	if err := db.First(&facility, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, facility)
}

func UpdateFacility(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var facility model.Facility
	if err := db.First(&facility, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	if err := c.ShouldBindJSON(&facility); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&facility)
	c.JSON(http.StatusOK, facility)
}

func DeleteFacility(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	if err := db.Delete(&model.Facility{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
