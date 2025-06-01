package model

import "time"

// Article 表示数据库中的 articles 表
type Article struct {
	ArticleID    int        `gorm:"column:article_id;primaryKey" json:"article_id"`
	Title        string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
	BodyText     string     `gorm:"column:body_text;type:text;not null" json:"body_text"`
	Category     string     `gorm:"column:category;type:varchar(100)" json:"category"`
	LikeCount    int        `gorm:"column:like_count;not null" json:"like_count"`
	ArticleImage []byte     `gorm:"column:article_image" json:"article_image,omitempty"`
	ImageFileID  *int       `gorm:"column:image_file_id" json:"image_file_id,omitempty"`
	CommentCount int        `gorm:"column:comment_count;not null" json:"comment_count"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`
	
	ImageURL     string     `gorm:"-" json:"image_url,omitempty"`
}

// ArticleReqCreate 文章创建请求
type ArticleReqCreate struct {
	Title        string `json:"title" binding:"required"`
	BodyText     string `json:"body_text" binding:"required"`
	Category     string `json:"category"`
	LikeCount    int    `json:"like_count"`
	ArticleImage []byte `json:"article_image"`
	ImageFileID  *int   `json:"image_file_id"`
	CommentCount int    `json:"comment_count"`
}

// ArticleReqEdit 文章更新请求
type ArticleReqEdit struct {
	ArticleID    int    `json:"article_id" binding:"required"`
	Title        string `json:"title"`
	BodyText     string `json:"body_text"`
	Category     string `json:"category"`
	LikeCount    int    `json:"like_count"`
	ArticleImage []byte `json:"article_image"`
	ImageFileID  *int   `json:"image_file_id"`
	CommentCount int    `json:"comment_count"`
}

// ArticleReqList 文章分页与搜索请求
type ArticleReqList struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Keyword  string `json:"keyword"`
}

// ArticleDetailRequest 获取单个文章请求
type ArticleDetailRequest struct {
	ArticleID int `json:"article_id" binding:"required"`
}

// ArticleCreateWithImageRequest 带图片上传的文章创建请求（multipart/form-data）
type ArticleCreateWithImageRequest struct {
	Title        string `form:"title" binding:"required"`
	BodyText     string `form:"body_text" binding:"required"`
	Category     string `form:"category"`
	LikeCount    int    `form:"like_count"`
	CommentCount int    `form:"comment_count"`
}
