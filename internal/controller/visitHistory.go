package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateVisitHistory godoc
// @Summary 新建访问记录
// @Description 新建一个访问记录
// @Tags VisitHistories
// @Accept json
// @Produce json
// @Param visit_history body model.VisitHistoryReqCreate true "访问记录信息"
// @Success 200 {object} model.Response[model.VisitHistory]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/visit_history [post]
func CreateVisitHistory(c *gin.Context) {
	var req model.VisitHistoryReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	history := model.VisitHistory{
		UserID:     req.UserID,
		FacilityID: req.FacilityID,
		ScanAt:     req.ScanAt,
		IsActive:   req.IsActive,
	}
	db := database.GetDB()
	if err := db.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.VisitHistory]{Success: true, Data: history})
}

// UpdateVisitHistory godoc
// @Summary 更新访问记录
// @Description 更新访问记录信息
// @Tags VisitHistories
// @Accept json
// @Produce json
// @Param visit_history body model.VisitHistoryReqEdit true "访问记录信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/visit_history [put]
func UpdateVisitHistory(c *gin.Context) {
	var req model.VisitHistoryReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var history model.VisitHistory
	if err := db.First(&history, req.HistoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "访问记录不存在"})
		return
	}
	db.Model(&history).Updates(map[string]interface{}{
		"user_id":     req.UserID,
		"facility_id": req.FacilityID,
		"scan_at":     req.ScanAt,
		"is_active":   req.IsActive,
	})
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteVisitHistory godoc
// @Summary 删除访问记录
// @Description 删除一个访问记录
// @Tags VisitHistories
// @Accept json
// @Produce json
// @Param history_id path int true "访问记录ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/visit_history/{history_id} [delete]
func DeleteVisitHistory(c *gin.Context) {
	id := c.Param("history_id")
	historyID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: "参数错误"})
		return
	}
	db := database.GetDB()
	if err := db.Delete(&model.VisitHistory{}, historyID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetVisitHistory godoc
// @Summary 获取访问记录
// @Description 获取单个访问记录
// @Tags VisitHistories
// @Accept json
// @Produce json
// @Param history_id path int true "访问记录ID"
// @Success 200 {object} model.Response[model.VisitHistory]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/visit_history/{history_id} [get]
func GetVisitHistory(c *gin.Context) {
	id := c.Param("history_id")
	historyID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: "参数错误"})
		return
	}
	db := database.GetDB()
	var history model.VisitHistory
	if err := db.First(&history, historyID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "访问记录不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.VisitHistory]{Success: true, Data: history})
}

// ListVisitHistories godoc
// @Summary 获取访问记录列表
// @Description 获取访问记录分页列表
// @Tags VisitHistories
// @Accept json
// @Produce json
// @Param req body model.VisitHistoryReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.VisitHistory]
// @Failure 400 {object} model.BaseResponse
// @Router /api/visit_history/list [post]
func ListVisitHistories(c *gin.Context) {
	var req model.VisitHistoryReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var histories []model.VisitHistory
	var total int64
	db.Model(&model.VisitHistory{}).Count(&total)
	db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&histories)

	c.JSON(http.StatusOK, model.ListResponse[model.VisitHistory]{
		Success: true,
		Total:   total,
		List:    histories,
	})
}
