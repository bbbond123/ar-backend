package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateMenu godoc
// @Summary 新建菜单
// @Description 新建一个菜单
// @Tags Menus
// @Accept json
// @Produce json
// @Param menu body model.MenuReqCreate true "菜单信息"
// @Success 200 {object} model.Response[model.Menu]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/menus [post]
func CreateMenu(c *gin.Context) {
	var req model.MenuReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	menu := model.Menu{
		MenuName:     req.MenuName,
		MenuCode:     req.MenuCode,
		DisplayOrder: req.DisplayOrder,
		IsActive:     req.IsActive,
	}
	db := database.GetDB()
	if err := db.Create(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Menu]{Success: true, Data: menu})
}

// UpdateMenu godoc
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags Menus
// @Accept json
// @Produce json
// @Param menu body model.MenuReqEdit true "菜单信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/menus [put]
func UpdateMenu(c *gin.Context) {
	var req model.MenuReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var menu model.Menu
	if err := db.First(&menu, req.MenuID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "菜单不存在"})
		return
	}
	db.Model(&menu).Updates(req)
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteMenu godoc
// @Summary 删除菜单
// @Description 删除一个菜单
// @Tags Menus
// @Accept json
// @Produce json
// @Param menu_id path int true "菜单ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/menus/{menu_id} [delete]
func DeleteMenu(c *gin.Context) {
	id := c.Param("menu_id")
	menuID, _ := strconv.Atoi(id)
	db := database.GetDB()
	if err := db.Delete(&model.Menu{}, menuID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetMenu godoc
// @Summary 获取菜单
// @Description 获取单个菜单信息
// @Tags Menus
// @Accept json
// @Produce json
// @Param menu_id path int true "菜单ID"
// @Success 200 {object} model.Response[model.Menu]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/menus/{menu_id} [get]
func GetMenu(c *gin.Context) {
	id := c.Param("menu_id")
	menuID, _ := strconv.Atoi(id)
	db := database.GetDB()
	var menu model.Menu
	if err := db.First(&menu, menuID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "菜单不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Menu]{Success: true, Data: menu})
}

// ListMenus godoc
// @Summary 获取菜单列表
// @Description 获取菜单分页列表
// @Tags Menus
// @Accept json
// @Produce json
// @Param req body model.MenuReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.Menu]
// @Failure 400 {object} model.BaseResponse
// @Router /api/menus/list [post]
func ListMenus(c *gin.Context) {
	var req model.MenuReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var menus []model.Menu
	var total int64

	db.Model(&model.Menu{}).Count(&total)
	db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&menus)

	c.JSON(http.StatusOK, model.ListResponse[model.Menu]{
		Success: true,
		Total:   total,
		List:    menus,
	})
}
