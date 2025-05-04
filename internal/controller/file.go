package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateFile godoc
// @Summary 新建文件
// @Description 新建一个文件
// @Tags Files
// @Accept json
// @Produce json
// @Param file body model.FileReqCreate true "文件信息"
// @Success 200 {object} model.Response[model.File]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/files [post]
func CreateFile(c *gin.Context) {
	var req model.FileReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	file := model.File{
		FileName:  req.FileName,
		FileType:  req.FileType,
		FileSize:  req.FileSize,
		FileData:  req.FileData,
		Location:  req.Location,
		RelatedID: req.RelatedID,
	}
	db := database.GetDB()
	if err := db.Create(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.File]{Success: true, Data: file})
}

// UpdateFile godoc
// @Summary 更新文件
// @Description 更新文件信息
// @Tags Files
// @Accept json
// @Produce json
// @Param file body model.FileReqEdit true "文件信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/files [put]
func UpdateFile(c *gin.Context) {
	var req model.FileReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var file model.File
	if err := db.First(&file, req.FileID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "文件不存在"})
		return
	}
	db.Model(&file).Updates(req)
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteFile godoc
// @Summary 删除文件
// @Description 删除一个文件
// @Tags Files
// @Accept json
// @Produce json
// @Param file_id path int true "文件ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/files/{file_id} [delete]
func DeleteFile(c *gin.Context) {
	id := c.Param("file_id")
	fileID, _ := strconv.Atoi(id)
	db := database.GetDB()
	if err := db.Delete(&model.File{}, fileID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetFile godoc
// @Summary 获取文件
// @Description 获取单个文件信息
// @Tags Files
// @Accept json
// @Produce json
// @Param file_id path int true "文件ID"
// @Success 200 {object} model.Response[model.File]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/files/{file_id} [get]
func GetFile(c *gin.Context) {
	id := c.Param("file_id")
	fileID, _ := strconv.Atoi(id)
	db := database.GetDB()
	var file model.File
	if err := db.First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "文件不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.File]{Success: true, Data: file})
}

// ListFiles godoc
// @Summary 获取文件列表
// @Description 获取文件分页列表
// @Tags Files
// @Accept json
// @Produce json
// @Param req body model.FileReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.File]
// @Failure 400 {object} model.BaseResponse
// @Router /api/files/list [post]
func ListFiles(c *gin.Context) {
	var req model.FileReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var files []model.File
	var total int64

	db.Model(&model.File{}).Count(&total)
	db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&files)

	c.JSON(http.StatusOK, model.ListResponse[model.File]{
		Success: true,
		Total:   total,
		List:    files,
	})
}
