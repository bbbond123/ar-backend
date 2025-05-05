package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary 新建用户
// @Description 新建一个用户
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.UserReqCreate true "用户信息"
// @Success 200 {object} model.Response[model.User]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/users [post]
func CreateUser(c *gin.Context) {
	var req model.UserReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	user := model.User{
		Name:        req.Name,
		NameKana:    req.NameKana,
		Address:     req.Address,
		Gender:      &req.Gender,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Password:    req.Password,
		GoogleID:    req.GoogleID,
		AppleID:     req.AppleID,
		Provider:    req.Provider,
		Status:      req.Status,
	}

	db := database.GetDB()
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.Response[model.User]{Success: true, Data: user})
}

// UpdateUser godoc
// @Summary 更新用户
// @Description 更新一个用户信息
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.UserReqEdit true "用户信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/users [put]
func UpdateUser(c *gin.Context) {
	var req model.UserReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	db := database.GetDB()
	var user model.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "用户不存在"})
		return
	}

	db.Model(&user).Updates(req)
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 删除一个用户
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/users/{user_id} [delete]
func DeleteUser(c *gin.Context) {
	id := c.Param("user_id")
	db := database.GetDB()
	if err := db.Delete(&model.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetUser godoc
// @Summary 获取用户信息
// @Description 获取单个用户信息
// @Tags Users
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Success 200 {object} model.Response[model.User]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/users/{user_id} [get]
func GetUser(c *gin.Context) {
	id := c.Param("user_id")
	db := database.GetDB()
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, model.Response[model.User]{Success: true, Data: user})
}

// ListUsers godoc
// @Summary 获取用户列表
// @Description 获取用户分页列表
// @Tags Users
// @Accept json
// @Produce json
// @Param req body model.UserReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.User]
// @Failure 400 {object} model.BaseResponse
// @Router /api/users/list [post]
func ListUsers(c *gin.Context) {
	var req model.UserReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	db := database.GetDB()
	var users []model.User
	var total int64

	query := db.Model(&model.User{})
	if req.Keyword != "" {
		query = query.Where("name ILIKE ?", "%"+req.Keyword+"%")
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&users)

	c.JSON(http.StatusOK, model.ListResponse[model.User]{
		Success: true,
		Total:   total,
		List:    users,
	})
}

// UserProfile godoc
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} model.Response[model.User]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/auth/user/profile [get]
func UserProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	db := database.GetDB()
	var user model.User
	db.First(&user, userID)
	c.JSON(http.StatusOK, model.Response[model.User]{Success: true, Data: user})

}
