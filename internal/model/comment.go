package model

import "time"

// Comment 表示数据库中的 comments 表
type Comment struct {
	CommentID        int        `gorm:"column:comment_id;primaryKey" json:"comment_id"`
	ArticleID        int        `gorm:"column:article_id;not null" json:"article_id"`
	UserID           int        `gorm:"column:user_id;not null" json:"user_id"`
	CommentText      string     `gorm:"column:comment_text;type:text;not null" json:"comment_text"`
	CreatedAt        time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"updated_at"`
	IsPublished      bool       `gorm:"column:is_published;not null" json:"is_published"`
	ReplyToCommentID *int       `gorm:"column:reply_to_comment_id" json:"reply_to_comment_id"`
}

// CommentReqCreate 新建评论请求
type CommentReqCreate struct {
	ArticleID        int    `json:"article_id" binding:"required"`
	UserID           int    `json:"user_id" binding:"required"`
	CommentText      string `json:"comment_text" binding:"required"`
	IsPublished      bool   `json:"is_published" binding:"required"`
	ReplyToCommentID *int   `json:"reply_to_comment_id"`
}

// CommentReqEdit 更新评论请求
type CommentReqEdit struct {
	CommentID        int    `json:"comment_id" binding:"required"`
	CommentText      string `json:"comment_text"`
	IsPublished      bool   `json:"is_published"`
	ReplyToCommentID *int   `json:"reply_to_comment_id"`
}

// CommentReqList 评论分页与搜索请求
type CommentReqList struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Keyword  string `json:"keyword"`
}

// CommentDetailRequest 获取单个评论请求
type CommentDetailRequest struct {
	CommentID int `json:"comment_id" binding:"required"`
}
