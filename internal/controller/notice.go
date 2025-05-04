package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateNotice godoc
// @Summary 新建通知
// @Description 新建一条通知
// @Tags Notices
// @Accept json
// @Produce json
// @Param notice body model.NoticeReqCreate true "通知信息"
// @Success 200 {object} model.Response[model.Notice]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/notices [post]
func CreateNotice(c *gin.Context) {
	var req model.NoticeReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	notice := model.Notice{
		Title:       req.Title,
		Content:     req.Content,
		NoticeType:  req.NoticeType,
		UserID:      req.UserID,
		PublishedAt: req.PublishedAt,
		IsActive:    req.IsActive,
		IsRead:      req.IsRead,
	}
	db := database.GetDB()
	if err := db.Create(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Notice]{Success: true, Data: notice})
}

// UpdateNotice godoc
// @Summary 更新通知
// @Description 更新通知内容
// @Tags Notices
// @Accept json
// @Produce json
// @Param notice body model.NoticeReqEdit true "通知信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/notices [put]
func UpdateNotice(c *gin.Context) {
	var req model.NoticeReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var notice model.Notice
	if err := db.First(&notice, req.NoticeID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "通知不存在"})
		return
	}
	db.Model(&notice).Updates(req)
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteNotice godoc
// @Summary 删除通知
// @Description 删除一条通知
// @Tags Notices
// @Accept json
// @Produce json
// @Param notice_id path int true "通知ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/notices/{notice_id} [delete]
func DeleteNotice(c *gin.Context) {
	id := c.Param("notice_id")
	noticeID, _ := strconv.Atoi(id)
	db := database.GetDB()
	if err := db.Delete(&model.Notice{}, noticeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetNotice godoc
// @Summary 获取通知
// @Description 获取单个通知信息
// @Tags Notices
// @Accept json
// @Produce json
// @Param notice_id path int true "通知ID"
// @Success 200 {object} model.Response[model.Notice]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/notices/{notice_id} [get]
func GetNotice(c *gin.Context) {
	id := c.Param("notice_id")
	noticeID, _ := strconv.Atoi(id)
	db := database.GetDB()
	var notice model.Notice
	if err := db.First(&notice, noticeID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "通知不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Notice]{Success: true, Data: notice})
}

// ListNotices godoc
// @Summary 获取通知列表
// @Description 获取通知分页列表
// @Tags Notices
// @Accept json
// @Produce json
// @Param req body model.NoticeReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.Notice]
// @Failure 400 {object} model.BaseResponse
// @Router /api/notices/list [post]
func ListNotices(c *gin.Context) {
	var req model.NoticeReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var notices []model.Notice
	var total int64

	db.Model(&model.Notice{}).Count(&total)
	db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&notices)

	c.JSON(http.StatusOK, model.ListResponse[model.Notice]{
		Success: true,
		Total:   total,
		List:    notices,
	})
}
