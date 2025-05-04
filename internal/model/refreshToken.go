package model

import "time"

// RefreshToken 表示数据库中的 refresh_tokens 表
type RefreshToken struct {
	TokenID      int       `gorm:"column:token_id;primaryKey" json:"token_id"`
	UserID       int       `gorm:"column:user_id;not null" json:"user_id"`
	RefreshToken string    `gorm:"column:refresh_token;type:varchar(255);not null" json:"refresh_token"`
	ExpiresAt    time.Time `gorm:"column:expires_at;not null" json:"expires_at"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	Revoked      bool      `gorm:"column:revoked;not null;default:false" json:"revoked"`
}

// RefreshTokenReqCreate 创建 Refresh Token 请求
type RefreshTokenReqCreate struct {
	UserID       int       `json:"user_id" binding:"required"`
	RefreshToken string    `json:"refresh_token" binding:"required"`
	ExpiresAt    time.Time `json:"expires_at" binding:"required"`
	Revoked      bool      `json:"revoked"`
}

// RefreshTokenReqEdit 更新 Refresh Token 请求
type RefreshTokenReqEdit struct {
	TokenID      int       `json:"token_id" binding:"required"`
	UserID       int       `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	Revoked      bool      `json:"revoked"`
}

// RefreshTokenReqList 分页与搜索请求
type RefreshTokenReqList struct {
	Page     int `json:"page" binding:"required"`
	PageSize int `json:"page_size" binding:"required"`
}
