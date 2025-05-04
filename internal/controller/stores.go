package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateStore godoc
// @Summary 新建商铺
// @Description 新建一个商铺
// @Tags Stores
// @Accept json
// @Produce json
// @Param store body model.StoreReqCreate true "商铺信息"
// @Success 200 {object} model.Response[model.Store]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/stores [post]
func CreateStore(c *gin.Context) {
	var req model.StoreReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	store := model.Store{
		StoreName:       req.StoreName,
		StoreCategory:   req.StoreCategory,
		Location:        req.Location,
		DescriptionText: req.Description,
		Address:         req.Address,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		BusinessHours:   req.BusinessHours,
		RatingScore:     req.RatingScore,
		PhoneNumber:     req.PhoneNumber,
	}
	db := database.GetDB()
	if err := db.Create(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response[model.Store]{Success: true, Data: store})
}

// UpdateStore godoc
// @Summary 更新商铺
// @Description 更新商铺信息
// @Tags Stores
// @Accept json
// @Produce json
// @Param store body model.StoreReqEdit true "商铺信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/stores [put]
func UpdateStore(c *gin.Context) {
	var req model.StoreReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var store model.Store
	if err := db.First(&store, req.StoreID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "商铺不存在"})
		return
	}
	db.Model(&store).Updates(req)
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteStore godoc
// @Summary 删除商铺
// @Description 删除一个商铺
// @Tags Stores
// @Accept json
// @Produce json
// @Param store_id path int true "商铺ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/stores/{store_id} [delete]
func DeleteStore(c *gin.Context) {
	id := c.Param("store_id")
	storeID, _ := strconv.Atoi(id)
	db := database.GetDB()
	if err := db.Delete(&model.Store{}, storeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetStore godoc
// @Summary 获取商铺信息
// @Description 获取单个商铺信息
// @Tags Stores
// @Accept json
// @Produce json
// @Param store_id path int true "商铺ID"
// @Success 200 {object} model.Response[model.Store]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/stores/{store_id} [get]
func GetStore(c *gin.Context) {
	id := c.Param("store_id")
	storeID, _ := strconv.Atoi(id)
	db := database.GetDB()
	var store model.Store
	if err := db.First(&store, storeID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "商铺不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Store]{Success: true, Data: store})
}

// ListStores godoc
// @Summary 获取商铺列表
// @Description 获取商铺分页列表
// @Tags Stores
// @Accept json
// @Produce json
// @Param req body model.StoreReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.Store]
// @Failure 400 {object} model.BaseResponse
// @Router /api/stores/list [post]
func ListStores(c *gin.Context) {
	var req model.StoreReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var stores []model.Store
	var total int64

	db.Model(&model.Store{}).Count(&total)
	db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&stores)

	c.JSON(http.StatusOK, model.ListResponse[model.Store]{
		Success: true,
		Total:   total,
		List:    stores,
	})
}

// @Router /api/stores/list [post]
func GetTagsByStore(c *gin.Context) {
	storeID := c.Param("storeID")

	var tags []model.Tag
	err := database.DB.Table("tags").
		Select("tags.*").
		Joins("JOIN taggings ON tags.tag_id = taggings.tag_id").
		Where("taggings.taggable_type = ? AND taggings.taggable_id = ?", "Store", storeID).
		Find(&tags).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"store_id": storeID,
		"tags":     tags,
	})
}

func parseUint(s string) uint {
	var id uint64
	id, _ = strconv.ParseUint(s, 10, 64)
	return uint(id)
}

// @Router /api/stores/list [post]
func AddTagToStore(c *gin.Context) {
	storeID := c.Param("storeID")

	var req struct {
		TagID uint `json:"tag_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	tagging := model.Tagging{
		TagID:        int(req.TagID),
		TaggableType: "Store",
		TaggableID:   parseUint(storeID),
	}

	if err := database.DB.Create(&tagging).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag added to store"})
}

func RemoveTagFromStore(c *gin.Context) {
	storeID := c.Param("storeID")
	tagID := c.Param("tagID")

	if err := database.DB.Where("taggable_type = ? AND taggable_id = ? AND tag_id = ?", "Store", storeID, tagID).
		Delete(&model.Tagging{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag removed from store"})
}
