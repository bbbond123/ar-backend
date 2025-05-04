package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateComment godoc
// @Summary 新建评论
// @Description 新建一条评论
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment body model.CommentReqCreate true "评论信息"
// @Success 200 {object} model.Response[model.Comment]
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/comments [post]
func CreateComment(c *gin.Context) {
	var req model.CommentReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	comment := model.Comment{
		ArticleID:        req.ArticleID,
		UserID:           req.UserID,
		CommentText:      req.CommentText,
		IsPublished:      req.IsPublished,
		ReplyToCommentID: req.ReplyToCommentID,
	}
	db := database.GetDB()
	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Comment]{Success: true, Data: comment})
}

// UpdateComment godoc
// @Summary 更新评论
// @Description 更新评论内容
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment body model.CommentReqEdit true "评论信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/comments [put]
func UpdateComment(c *gin.Context) {
	var req model.CommentReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var comment model.Comment
	if err := db.First(&comment, req.CommentID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "评论不存在"})
		return
	}
	db.Model(&comment).Updates(req)
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteComment godoc
// @Summary 删除评论
// @Description 删除一条评论
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment_id path int true "评论ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Router /api/comments/{comment_id} [delete]
func DeleteComment(c *gin.Context) {
	id := c.Param("comment_id")
	commentID, _ := strconv.Atoi(id)
	db := database.GetDB()
	if err := db.Delete(&model.Comment{}, commentID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetComment godoc
// @Summary 获取单个评论
// @Description 获取一条评论信息
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment_id path int true "评论ID"
// @Success 200 {object} model.Response[model.Comment]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/comments/{comment_id} [get]
func GetComment(c *gin.Context) {
	id := c.Param("comment_id")
	commentID, _ := strconv.Atoi(id)
	db := database.GetDB()
	var comment model.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "评论不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Comment]{Success: true, Data: comment})
}

// ListComments godoc
// @Summary 获取评论列表
// @Description 获取评论分页列表
// @Tags Comments
// @Accept json
// @Produce json
// @Param req body model.CommentReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.Comment]
// @Failure 400 {object} model.BaseResponse
// @Router /api/comments/list [post]
func ListComments(c *gin.Context) {
	var req model.CommentReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var comments []model.Comment
	var total int64

	db.Model(&model.Comment{}).Count(&total)
	db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&comments)

	c.JSON(http.StatusOK, model.ListResponse[model.Comment]{
		Success: true,
		Total:   total,
		List:    comments,
	})
}
