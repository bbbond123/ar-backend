package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateFacility godoc
// @Summary 新建设施
// @Description 新增一个设施记录
// @Tags Facilities
// @Accept json
// @Produce json
// @Param facility body model.FacilityCreateRequest true "设施信息"
// @Success 200 {object} model.Response[model.Facility]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/facilities [post]
func CreateFacility(c *gin.Context) {
	var req model.FacilityCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	db := database.GetDB()
	facility := model.Facility{
		FacilityName:    req.FacilityName,
		Location:        req.Location,
		DescriptionText: req.DescriptionText,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		PersonID:        req.PersonID,
	}

	if err := db.Create(&facility).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Facility]{Success: true, Data: facility})
}

// UpdateFacility godoc
// @Summary 更新设施
// @Description 更新指定ID的设施
// @Tags Facilities
// @Accept json
// @Produce json
// @Param id path int true "设施ID"
// @Param facility body model.FacilityUpdateRequest true "设施信息"
// @Success 200 {object} model.Response[model.Facility]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/facilities/{id} [put]
func UpdateFacility(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var facility model.Facility
	if err := db.First(&facility, id).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "Not found"})
		return
	}

	var req model.FacilityUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	facility.FacilityName = req.FacilityName
	facility.Location = req.Location
	facility.DescriptionText = req.DescriptionText
	facility.Latitude = req.Latitude
	facility.Longitude = req.Longitude
	facility.PersonID = req.PersonID

	db.Save(&facility)
	c.JSON(http.StatusOK, model.Response[model.Facility]{Success: true, Data: facility})
}

// DeleteFacility godoc
// @Summary 删除设施
// @Description 删除指定ID的设施
// @Tags Facilities
// @Accept json
// @Produce json
// @Param id path int true "设施ID"
// @Success 200 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/facilities/{id} [delete]
func DeleteFacility(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	if err := db.Delete(&model.Facility{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.BaseResponse{Success: true, ErrMessage: "Deleted"})
}

// GetFacility godoc
// @Summary 获取单个设施
// @Description 通过ID获取设施详情
// @Tags Facilities
// @Accept json
// @Produce json
// @Param id path int true "设施ID"
// @Success 200 {object} model.Response[model.Facility]
// @Failure 404 {object} model.BaseResponse
// @Router /api/facilities/{id} [get]
func GetFacility(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	var facility model.Facility
	if err := db.First(&facility, id).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "Not found"})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Facility]{Success: true, Data: facility})
}

// ListFacilities godoc
// @Summary 获取设施列表
// @Description 获取设施列表（分页 + 搜索）
// @Tags Facilities
// @Accept json
// @Produce json
// @Param query body model.FacilityQueryRequest true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.Facility]
// @Failure 400 {object} model.BaseResponse
// @Router /api/facilities/list [post]
func ListFacilities(c *gin.Context) {
	var query model.FacilityQueryRequest
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var facilities []model.Facility
	var total int64

	db.Model(&model.Facility{}).
		Where("facility_name LIKE ?", "%"+query.Keyword+"%").
		Count(&total).
		Limit(query.PageSize).
		Offset((query.Page - 1) * query.PageSize).
		Find(&facilities)

	c.JSON(http.StatusOK, model.ListResponse[model.Facility]{
		Total:   total,
		List:    facilities,
		Success: true,
	})
}
