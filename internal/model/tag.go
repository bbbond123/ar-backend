package model

import "time"

// Tag 表示 tags 表
type Tag struct {
	TagID     int        `gorm:"column:tag_id;primaryKey" json:"tag_id"`
	TagName   string     `gorm:"column:tag_name;type:varchar(50);not null" json:"tag_name"`
	IsActive  bool       `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TagReqCreate 创建Tag请求
type TagReqCreate struct {
	TagName  string `json:"tag_name" binding:"required"`
	IsActive bool   `json:"is_active"`
}

// TagReqEdit 更新Tag请求
type TagReqEdit struct {
	TagID    int    `json:"tag_id" binding:"required"`
	TagName  string `json:"tag_name"`
	IsActive bool   `json:"is_active"`
}

// TagReqList Tag分页请求
type TagReqList struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Keyword  string `json:"keyword"`
}

// Tagging 表示 taggings 表
type Tagging struct {
	TaggingID    int        `gorm:"column:tagging_id;primaryKey" json:"tagging_id"`
	TagID        int        `gorm:"column:tag_id;not null" json:"tag_id"`
	TaggableType string     `gorm:"column:taggable_type;type:varchar(50);not null" json:"taggable_type"`
	TaggableID   int        `gorm:"column:taggable_id;not null" json:"taggable_id"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TaggingReqCreate 创建Tagging请求
type TaggingReqCreate struct {
	TagID        int    `json:"tag_id" binding:"required"`
	TaggableType string `json:"taggable_type" binding:"required"`
	TaggableID   int    `json:"taggable_id" binding:"required"`
}

// TaggingReqEdit 更新Tagging请求
type TaggingReqEdit struct {
	TaggingID    int    `json:"tagging_id" binding:"required"`
	TagID        int    `json:"tag_id"`
	TaggableType string `json:"taggable_type"`
	TaggableID   int    `json:"taggable_id"`
}

// TaggingReqList Tagging分页请求
type TaggingReqList struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Keyword  string `json:"keyword"`
}
