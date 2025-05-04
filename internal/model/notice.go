package model

import "time"

// Notice 表示数据库中的 notices 表
type Notice struct {
	NoticeID    int        `gorm:"column:notice_id;primaryKey" json:"notice_id"`
	Title       string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Content     string     `gorm:"column:content;type:text;not null" json:"content"`
	NoticeType  bool       `gorm:"column:notice_type;not null" json:"notice_type"`
	UserID      *int       `gorm:"column:user_id" json:"user_id"`
	PublishedAt time.Time  `gorm:"column:published_at;not null" json:"published_at"`
	IsActive    bool       `gorm:"column:is_active;not null" json:"is_active"`
	IsRead      bool       `gorm:"column:is_read;default:false" json:"is_read"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// NoticeReqCreate 新建请求
type NoticeReqCreate struct {
	Title       string    `json:"title" binding:"required"`
	Content     string    `json:"content" binding:"required"`
	NoticeType  bool      `json:"notice_type" binding:"required"`
	UserID      *int      `json:"user_id"`
	PublishedAt time.Time `json:"published_at" binding:"required"`
	IsActive    bool      `json:"is_active" binding:"required"`
	IsRead      bool      `json:"is_read"`
}

// NoticeReqEdit 更新请求
type NoticeReqEdit struct {
	NoticeID    int       `json:"notice_id" binding:"required"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	NoticeType  bool      `json:"notice_type"`
	UserID      *int      `json:"user_id"`
	PublishedAt time.Time `json:"published_at"`
	IsActive    bool      `json:"is_active"`
	IsRead      bool      `json:"is_read"`
}

// NoticeReqList 分页与搜索请求
type NoticeReqList struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Keyword  string `json:"keyword"`
}

// NoticeDetailRequest 单个查询请求
type NoticeDetailRequest struct {
	NoticeID int `json:"notice_id" binding:"required"`
}
