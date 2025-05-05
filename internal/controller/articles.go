package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/database"
	"net/http"
	"strconv"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// CreateArticle godoc
// @Summary 新建文章
// @Description 新建一个文章
// @Tags Articles
// @Accept json
// @Produce json
// @Param article body model.ArticleReqCreate true "文章信息"
// @Success 200 {object} model.Response[model.Article]
// @Failure 400 {object} model.BaseResponse
// @Failure 401 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Security ApiKeyAuth
// @Router /api/articles [post]
func CreateArticle(c *gin.Context) {
	// 1. 获取并校验access token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, model.BaseResponse{Success: false, ErrMessage: "未登录，缺少token"})
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	type UserIDClaims struct {
		UserID int `json:"user_id"`
		jwt.RegisteredClaims
	}
	claims := &UserIDClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		c.JSON(401, model.BaseResponse{Success: false, ErrMessage: "token无效或已过期"})
		return
	}
	// 2. 解析请求体
	var req model.ArticleReqCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	// 3. 创建文章
	article := model.Article{
		Title:        req.Title,
		BodyText:     req.BodyText,
		Category:     req.Category,
		LikeCount:    req.LikeCount,
		ArticleImage: req.ArticleImage,
		CommentCount: req.CommentCount,
		// 可选：UserID: claims.UserID,
	}
	db := database.GetDB()
	if err := db.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Article]{Success: true, Data: article})
}

// UpdateArticle godoc
// @Summary 更新文章
// @Description 更新文章信息
// @Tags Articles
// @Accept json
// @Produce json
// @Param article body model.ArticleReqEdit true "文章信息"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 401 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Security ApiKeyAuth
// @Router /api/articles [put]
func UpdateArticle(c *gin.Context) {
	var req model.ArticleReqEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var article model.Article
	if err := db.First(&article, req.ArticleID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "文章不存在"})
		return
	}
	db.Model(&article).Updates(req)
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// DeleteArticle godoc
// @Summary 删除文章
// @Description 删除一个文章
// @Tags Articles
// @Accept json
// @Produce json
// @Param article_id path int true "文章ID"
// @Success 200 {object} model.BaseResponse
// @Failure 400 {object} model.BaseResponse
// @Failure 401 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Security ApiKeyAuth
// @Router /api/articles/{article_id} [delete]
func DeleteArticle(c *gin.Context) {
	id := c.Param("article_id")
	articleID, _ := strconv.Atoi(id)
	db := database.GetDB()
	if err := db.Delete(&model.Article{}, articleID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.BaseResponse{Success: true})
}

// GetArticle godoc
// @Summary 获取文章信息
// @Description 获取单个文章信息
// @Tags Articles
// @Accept json
// @Produce json
// @Param article_id path int true "文章ID"
// @Success 200 {object} model.Response[model.Article]
// @Failure 400 {object} model.BaseResponse
// @Failure 404 {object} model.BaseResponse
// @Router /api/articles/{article_id} [get]
func GetArticle(c *gin.Context) {
	id := c.Param("article_id")
	articleID, _ := strconv.Atoi(id)
	db := database.GetDB()
	var article model.Article
	if err := db.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "文章不存在"})
		return
	}
	c.JSON(http.StatusOK, model.Response[model.Article]{Success: true, Data: article})
}

// ListArticles godoc
// @Summary 获取文章列表
// @Description 获取文章分页列表
// @Tags Articles
// @Accept json
// @Produce json
// @Param req body model.ArticleReqList true "分页与搜索"
// @Success 200 {object} model.ListResponse[model.Article]
// @Failure 400 {object} model.BaseResponse
// @Router /api/articles/list [post]
func ListArticles(c *gin.Context) {
	var req model.ArticleReqList
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}
	db := database.GetDB()
	var articles []model.Article
	var total int64

	db.Model(&model.Article{}).Count(&total)
	db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&articles)

	c.JSON(http.StatusOK, model.ListResponse[model.Article]{
		Success: true,
		Total:   total,
		List:    articles,
	})
}
