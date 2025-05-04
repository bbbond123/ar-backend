package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateRefreshToken godoc
// @Summary 新建Refresh Token
// @Description 新建一个Refresh Token
// @Tags RefreshTokens
// @Accept json
// @Produce json
// @Param refresh_token body model.RefreshTokenReqCreate true "Refresh Token信息"
// @Success 200 {object} model.Response[model.RefreshToken]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/refresh_tokens [post]
func CreateRefreshToken(c *gin.Context) {
	var req model.RefreshTokenReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	token := model.RefreshToken{
		UserID:       req.UserID,
		RefreshToken: req.RefreshToken,
		ExpiresAt:    req.ExpiresAt,
		Revoked:      req.Revoked,
	}
	db := database.GetDB()
	if err := db.Create(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.RefreshToken]{Success: true, Data: token})
}

// UpdateRefreshToken godoc
// @Summary 更新Refresh Token
// @Description 更新Refresh Token
// @Tags RefreshTokens
// @Accept json
// @Produce json
// @Param refresh_token body model.RefreshTokenReqEdit true "Refresh Token信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/refresh_tokens [put]
func UpdateRefreshToken(c *gin.Context) {
	var req model.RefreshTokenReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var token model.RefreshToken
	if err := db.First(&token, req.TokenID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "Token不存在"})
		return
	}
	db.Model(&token).Updates(req)
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteRefreshToken godoc
// @Summary 删除Refresh Token
// @Description 删除一个Refresh Token
// @Tags RefreshTokens
// @Accept json
// @Produce json
// @Param token_id path int true "Token ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/refresh_tokens/{token_id} [delete]
func DeleteRefreshToken(c *gin.Context) {
	id := c.Param("token_id")
	tokenID, _ := strconv.Atoi(id)
	db := database.GetDB()
	if err := db.Delete(&model.RefreshToken{}, tokenID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetRefreshToken godoc
// @Summary 获取Refresh Token
// @Description 获取单个Refresh Token
// @Tags RefreshTokens
// @Accept json
// @Produce json
// @Param token_id path int true "Token ID"
// @Success 200 {object} model.Response[model.RefreshToken]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/refresh_tokens/{token_id} [get]
func GetRefreshToken(c *gin.Context) {
	id := c.Param("token_id")
	tokenID, _ := strconv.Atoi(id)
	db := database.GetDB()
	var token model.RefreshToken
	if err := db.First(&token, tokenID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "Token不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.RefreshToken]{Success: true, Data: token})
}

// ListRefreshTokens godoc
// @Summary 获取Refresh Token列表
// @Description 获取Refresh Token分页列表
// @Tags RefreshTokens
// @Accept json
// @Produce json
// @Param req body model.RefreshTokenReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.RefreshToken]
// @Failure 400 {object} model.BaseResponse
// @Router /api/refresh_tokens/list [post]
func ListRefreshTokens(c *gin.Context) {
	var req model.RefreshTokenReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var tokens []model.RefreshToken
	var total int64

	db.Model(&model.RefreshToken{}).Count(&total)
	db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&tokens)

	c.JSON(http.StatusOK, model.ListResponse[model.RefreshToken]{
		Success: true,
		Total:   total,
		List:    tokens,
	})
}
