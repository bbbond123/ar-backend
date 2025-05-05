package model

import "time"

// Language 表示 languages 表
type Language struct {
	LanguageID   int        `gorm:"column:language_id;primaryKey" json:"language_id"`
	LanguageName string     `gorm:"column:language_name;type:varchar(50);not null" json:"language_name"`
	DisplayOrder *int       `gorm:"column:display_order" json:"display_order"`
	IsActive     bool       `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// LanguageReqCreate 创建语言请求
type LanguageReqCreate struct {
	LanguageName string `json:"language_name" binding:"required"`
	DisplayOrder *int   `json:"display_order"`
	IsActive     bool   `json:"is_active"`
}

// LanguageReqEdit 更新语言请求
type LanguageReqEdit struct {
	LanguageID   int    `json:"language_id" binding:"required"`
	LanguageName string `json:"language_name"`
	DisplayOrder *int   `json:"display_order"`
	IsActive     bool   `json:"is_active"`
}

// LanguageReqList 语言分页请求
type LanguageReqList struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Keyword  string `json:"keyword"`
}
