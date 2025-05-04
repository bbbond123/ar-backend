package model

import "time"

// User 表示数据库中的 users 表
type User struct {
	UserID      int        `gorm:"column:user_id;primaryKey" json:"user_id"`
	Name        string     `gorm:"column:name" json:"name"`
	NameKana    string     `gorm:"column:name_kana" json:"name_kana"`
	Birth       *time.Time `gorm:"column:birth" json:"birth"`
	Address     string     `gorm:"column:address" json:"address"`
	Gender      *string    `gorm:"column:gender" json:"gender"`
	PhoneNumber string     `gorm:"column:phone_number" json:"phone_number"`
	Email       string     `gorm:"column:email;not null;unique" json:"email"`
	Password    string     `gorm:"column:password" json:"password"`
	GoogleID    string     `gorm:"column:google_id" json:"google_id"`
	AppleID     string     `gorm:"column:apple_id" json:"apple_id"`
	Provider    string     `gorm:"column:provider;not null" json:"provider"`
	Status      string     `gorm:"column:status;not null" json:"status"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// UserReqCreate 用户创建请求
type UserReqCreate struct {
	Name        string `json:"name"`
	NameKana    string `json:"name_kana"`
	Birth       string `json:"birth"`
	Address     string `json:"address"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password"`
	GoogleID    string `json:"google_id"`
	AppleID     string `json:"apple_id"`
	Provider    string `json:"provider" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

// UserReqEdit 用户更新请求
type UserReqEdit struct {
	UserID      int    `json:"user_id" binding:"required"`
	Name        string `json:"name"`
	NameKana    string `json:"name_kana"`
	Birth       string `json:"birth"`
	Address     string `json:"address"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	GoogleID    string `json:"google_id"`
	AppleID     string `json:"apple_id"`
	Provider    string `json:"provider"`
	Status      string `json:"status"`
}

// UserReqList 用户分页与搜索请求
type UserReqList struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Keyword  string `json:"keyword"`
}

// UserDetailRequest 获取单个用户请求
type UserDetailRequest struct {
	UserID int `json:"user_id" binding:"required"`
}
