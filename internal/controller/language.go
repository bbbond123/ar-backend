package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateLanguage godoc
// @Summary 创建语言
// @Description 创建一个新语言
// @Tags Languages
// @Accept json
// @Produce json
// @Param language body model.LanguageReqCreate true "语言信息"
// @Success 200 {object} model.Response[model.Language]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/languages [post]
func CreateLanguage(c *gin.Context) {
	var req model.LanguageReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	language := model.Language{
		LanguageName: req.LanguageName,
		DisplayOrder: req.DisplayOrder,
		IsActive:     req.IsActive,
	}
	db := database.GetDB()
	if err := db.Create(&language).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Language]{Success: true, Data: language})
}

// UpdateLanguage godoc
// @Summary 更新语言
// @Description 更新语言信息
// @Tags Languages
// @Accept json
// @Produce json
// @Param language body model.LanguageReqEdit true "语言信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/languages [put]
func UpdateLanguage(c *gin.Context) {
	var req model.LanguageReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	db := database.GetDB()
	var language model.Language
	if err := db.First(&language, req.LanguageID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "语言不存在"})
		return
	}

	db.Model(&language).Updates(model.Language{
		LanguageName: req.LanguageName,
		DisplayOrder: req.DisplayOrder,
		IsActive:     req.IsActive,
	})

	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteLanguage godoc
// @Summary 删除语言
// @Description 删除一个语言
// @Tags Languages
// @Accept json
// @Produce json
// @Param language_id path int true "语言ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/languages/{language_id} [delete]
func DeleteLanguage(c *gin.Context) {
	id := c.Param("language_id")
	languageID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: "参数错误"})
		return
	}
	db := database.GetDB()
	if err := db.Delete(&model.Language{}, languageID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetLanguage godoc
// @Summary 获取语言
// @Description 获取单个语言
// @Tags Languages
// @Accept json
// @Produce json
// @Param language_id path int true "语言ID"
// @Success 200 {object} model.Response[model.Language]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/languages/{language_id} [get]
func GetLanguage(c *gin.Context) {
	id := c.Param("language_id")
	languageID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: "参数错误"})
		return
	}
	db := database.GetDB()
	var language model.Language
	if err := db.First(&language, languageID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "语言不存在"})
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Language]{Success: true, Data: language})
}

// ListLanguages godoc
// @Summary 获取语言列表
// @Description 获取语言分页列表
// @Tags Languages
// @Accept json
// @Produce json
// @Param req body model.LanguageReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.Language]
// @Failure 400 {object} model.BaseResponse
// @Router /api/languages/list [post]
func ListLanguages(c *gin.Context) {
	var req model.LanguageReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var languages []model.Language
	var total int64
	db.Model(&model.Language{}).Count(&total)
	db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&languages)

	c.JSON(http.StatusOK, model.ListResponse[model.Language]{
		Success: true,
		Total:   total,
		List:    languages,
	})
}
