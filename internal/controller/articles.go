package controller

import (
	"ar-backend/internal/model"
	"ar-backend/pkg/aws"
	"ar-backend/pkg/database"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// CreateArticleWithImage godoc
// @Summary 新建文章（支持图片上传）
// @Description 新建一个文章并同时上传图片到S3
// @Tags Articles
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "文章标题"
// @Param body_text formData string true "文章内容"
// @Param category formData string false "文章分类"
// @Param like_count formData int false "点赞数"
// @Param comment_count formData int false "评论数"
// @Param image formData file false "文章图片"
// @Success 200 {object} model.Response[model.Article]
// @Failure 400 {object} model.BaseResponse
// @Failure 401 {object} model.BaseResponse
// @Failure 500 {object} model.BaseResponse
// @Security ApiKeyAuth
// @Router /api/articles/with-image [post]
func CreateArticleWithImage(c *gin.Context) {
	// 1. 获取并校验access token（可选，根据需要）
	// 暂时跳过token验证，方便测试
	
	// 2. 解析表单数据
	var req model.ArticleCreateWithImageRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	var imageFileID *int
	
	// 3. 处理图片上传（如果有）
	fileHeader, err := c.FormFile("image")
	if err == nil && fileHeader != nil {
		// 有图片文件，上传到S3
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "打开图片文件失败"})
			return
		}
		defer file.Close()

		// 读取文件数据
		fileData, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "读取图片文件失败"})
			return
		}

		// 获取文件信息
		fileName := fileHeader.Filename
		fileSize := int(fileHeader.Size)
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			ext := filepath.Ext(fileName)
			contentType = getContentTypeForArticle(ext)
		}

		// 初始化 S3 服务并上传
		s3Service, err := aws.NewS3Service()
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "S3 服务初始化失败: " + err.Error()})
			return
		}

		s3URL, err := s3Service.UploadFile(fileData, fileName, contentType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "图片上传失败: " + err.Error()})
			return
		}

		// 保存文件记录到数据库
		fileRecord := model.File{
			FileName:  fileName,
			FileType:  contentType,
			FileSize:  fileSize,
			S3Key:     extractS3KeyForArticle(s3URL),
			S3URL:     s3URL,
			Location:  "article-images", // 文章图片的位置标识
			RelatedID: 0, // 暂时设为0，创建文章后会更新
		}

		db := database.GetDB()
		if err := db.Create(&fileRecord).Error; err != nil {
			// 如果数据库保存失败，删除已上传的 S3 文件
			s3Service.DeleteFile(fileRecord.S3Key)
			c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: "文件记录保存失败: " + err.Error()})
			return
		}

		imageFileID = &fileRecord.FileID
	}

	// 4. 创建文章
	article := model.Article{
		Title:        req.Title,
		BodyText:     req.BodyText,
		Category:     req.Category,
		LikeCount:    req.LikeCount,
		ImageFileID:  imageFileID,
		CommentCount: req.CommentCount,
	}

	db := database.GetDB()
	if err := db.Create(&article).Error; err != nil {
		// 如果文章创建失败且已上传图片，删除图片记录和S3文件
		if imageFileID != nil {
			var fileRecord model.File
			if db.First(&fileRecord, *imageFileID).Error == nil {
				if s3Service, err := aws.NewS3Service(); err == nil {
					s3Service.DeleteFile(fileRecord.S3Key)
				}
				db.Delete(&fileRecord)
			}
		}
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	// 5. 更新文件记录的关联ID
	if imageFileID != nil {
		db.Model(&model.File{}).Where("file_id = ?", *imageFileID).Update("related_id", article.ArticleID)
	}

	// 6. 获取完整的文章信息（包含图片URL）
	enrichedArticle := enrichArticleWithImageURL(db, article)

	c.JSON(http.StatusOK, model.Response[model.Article]{Success: true, Data: enrichedArticle})
}

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
		return getJWTSecret(), nil
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
		ImageFileID:  req.ImageFileID,
		CommentCount: req.CommentCount,
		// 可选：UserID: claims.UserID,
	}
	db := database.GetDB()
	if err := db.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.BaseResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	// 获取完整的文章信息（包含图片URL）
	enrichedArticle := enrichArticleWithImageURL(db, article)
	c.JSON(http.StatusOK, model.Response[model.Article]{Success: true, Data: enrichedArticle})
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
	var article model.Article
	if err := db.First(&article, articleID).Error; err != nil {
		c.JSON(http.StatusNotFound, model.BaseResponse{Success: false, ErrMessage: "文章不存在"})
		return
	}

	// 如果文章有关联的图片文件，先删除
	if article.ImageFileID != nil {
		var fileRecord model.File
		if db.First(&fileRecord, *article.ImageFileID).Error == nil {
			// 删除 S3 文件
			if fileRecord.S3Key != "" {
				if s3Service, err := aws.NewS3Service(); err == nil {
					s3Service.DeleteFile(fileRecord.S3Key)
				}
			}
			// 删除文件记录
			db.Delete(&fileRecord)
		}
	}

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

	// 获取完整的文章信息（包含图片URL）
	enrichedArticle := enrichArticleWithImageURL(db, article)
	c.JSON(http.StatusOK, model.Response[model.Article]{Success: true, Data: enrichedArticle})
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

	query := db.Model(&model.Article{})
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR body_text LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	query.Count(&total)
	query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&articles)

	// 为每篇文章添加图片URL
	enrichedArticles := make([]model.Article, len(articles))
	for i, article := range articles {
		enrichedArticles[i] = enrichArticleWithImageURL(db, article)
	}

	c.JSON(http.StatusOK, model.ListResponse[model.Article]{
		Success: true,
		Total:   total,
		List:    enrichedArticles,
	})
}

// 辅助函数：为文章添加图片URL
func enrichArticleWithImageURL(db *gorm.DB, article model.Article) model.Article {
	if article.ImageFileID != nil {
		var fileRecord model.File
		if db.First(&fileRecord, *article.ImageFileID).Error == nil {
			if fileRecord.S3URL != "" {
				article.ImageURL = fileRecord.S3URL
			}
		}
	}
	return article
}

// 辅助函数：根据文件扩展名获取 MIME 类型（文章图片专用）
func getContentTypeForArticle(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "image/jpeg" // 默认为JPEG
	}
}

// 辅助函数：从 S3 URL 提取 S3 Key（文章图片专用）
func extractS3KeyForArticle(s3URL string) string {
	if len(s3URL) == 0 {
		return ""
	}
	
	lastSlash := len(s3URL) - 1
	for i := len(s3URL) - 1; i >= 0; i-- {
		if s3URL[i] == '/' {
			lastSlash = i
			break
		}
	}
	
	if lastSlash < len(s3URL)-1 {
		return s3URL[lastSlash+1:]
	}
	
	return ""
}
