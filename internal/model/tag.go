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
