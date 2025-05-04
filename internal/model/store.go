package model

import "time"

// Store 表示数据库中的 stores 表
type Store struct {
	StoreID         int       `gorm:"column:store_id;primaryKey" json:"store_id"`
	StoreName       string    `gorm:"column:store_name;type:varchar(255);not null" json:"store_name"`
	StoreCategory   string    `gorm:"column:store_category;type:varchar(100);not null" json:"store_category"`
	Location        string    `gorm:"column:location;type:varchar(255);not null" json:"location"`
	DescriptionText string    `gorm:"column:description_text;type:text" json:"description"`
	Address         string    `gorm:"column:address;type:varchar(255);not null" json:"address"`
	Latitude        float64   `gorm:"column:latitude;type:decimal(10,6);not null" json:"latitude"`
	Longitude       float64   `gorm:"column:longitude;type:decimal(10,6);not null" json:"longitude"`
	BusinessHours   string    `gorm:"column:business_hours;type:varchar(100);not null" json:"business_hours"`
	RatingScore     float64   `gorm:"column:rating_score;type:decimal(3,2);not null" json:"rating_score"`
	PhoneNumber     string    `gorm:"column:phone_number;type:varchar(20);not null" json:"phone_number"`
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// StoreReqCreate 创建请求
type StoreReqCreate struct {
	StoreName     string  `json:"store_name" binding:"required"`
	StoreCategory string  `json:"store_category" binding:"required"`
	Location      string  `json:"location" binding:"required"`
	Description   string  `json:"description"`
	Address       string  `json:"address" binding:"required"`
	Latitude      float64 `json:"latitude" binding:"required"`
	Longitude     float64 `json:"longitude" binding:"required"`
	BusinessHours string  `json:"business_hours" binding:"required"`
	RatingScore   float64 `json:"rating_score" binding:"required"`
	PhoneNumber   string  `json:"phone_number" binding:"required"`
}

// StoreReqEdit 更新请求
type StoreReqEdit struct {
	StoreID       int     `json:"store_id" binding:"required"`
	StoreName     string  `json:"store_name"`
	StoreCategory string  `json:"store_category"`
	Location      string  `json:"location"`
	Description   string  `json:"description"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	BusinessHours string  `json:"business_hours"`
	RatingScore   float64 `json:"rating_score"`
	PhoneNumber   string  `json:"phone_number"`
}

// StoreReqList 查询请求
type StoreReqList struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Keyword  string `json:"keyword"`
}

// StoreDetailRequest 单个查询请求
type StoreDetailRequest struct {
	StoreID int `json:"store_id" binding:"required"`
}
