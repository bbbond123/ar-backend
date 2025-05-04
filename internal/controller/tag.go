package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ar-backend/internal/model"
	"ar-backend/pkg/database"

	"strconv"
)

// CreateTag godoc
// @Summary 创建标签
// @Description 创建一个新标签
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag body model.TagReqCreate true "标签信息"
// @Success 200 {object} model.Response[model.Tag]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/tags [post]
func CreateTag(c *gin.Context) {
	var req model.TagReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	tag := model.Tag{
		TagName:  req.TagName,
		IsActive: req.IsActive,
	}
	db := database.GetDB()
	if err := db.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Tag]{Success: true, Data: tag})
}

// UpdateTag godoc
// @Summary 更新标签
// @Description 更新标签信息
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag body model.TagReqEdit true "标签信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/tags [put]
func UpdateTag(c *gin.Context) {
	var req model.TagReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var tag model.Tag
	if err := db.First(&tag, req.TagID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "标签不存在"})
		return
	}
	db.Model(&tag).Updates(model.Tag{TagName: req.TagName, IsActive: req.IsActive})
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteTag godoc
// @Summary 删除标签
// @Description 删除一个标签
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag_id path int true "标签ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/tags/{tag_id} [delete]
func DeleteTag(c *gin.Context) {
	id := c.Param("tag_id")
	tagID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: "参数错误"})
		return
	}
	db := database.GetDB()
	if err := db.Delete(&model.Tag{}, tagID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetTag godoc
// @Summary 获取标签
// @Description 获取单个标签
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag_id path int true "标签ID"
// @Success 200 {object} model.Response[model.Tag]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/tags/{tag_id} [get]
func GetTag(c *gin.Context) {
	id := c.Param("tag_id")
	tagID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: "参数错误"})
		return
	}
	db := database.GetDB()
	var tag model.Tag
	if err := db.First(&tag, tagID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "标签不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Tag]{Success: true, Data: tag})
}

// ListTags godoc
// @Summary 获取标签列表
// @Description 获取标签分页列表
// @Tags Tags
// @Accept json
// @Produce json
// @Param req body model.TagReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.Tag]
// @Failure 400 {object} model.BaseResponse
// @Router /api/tags/list [post]
func ListTags(c *gin.Context) {
	var req model.TagReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var tags []model.Tag
	var total int64
	db.Model(&model.Tag{}).Count(&total)
	db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&tags)

	c.JSON(http.StatusOK, model.ListResponse[model.Tag]{
		Success: true,
		Total:   total,
		List:    tags,
	})
}
