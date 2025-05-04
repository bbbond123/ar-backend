package model

import "time"

// VisitHistory 表示数据库中的 visit_history 表
type VisitHistory struct {
	HistoryID  int        `gorm:"column:history_id;primaryKey" json:"history_id"`
	UserID     int        `gorm:"column:user_id;not null" json:"user_id"`
	FacilityID int        `gorm:"column:facility_id;not null" json:"facility_id"`
	ScanAt     time.Time  `gorm:"column:scan_at;not null" json:"scan_at"`
	IsActive   bool       `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt  time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// VisitHistoryReqCreate 新建访问记录请求
type VisitHistoryReqCreate struct {
	UserID     int       `json:"user_id" binding:"required"`
	FacilityID int       `json:"facility_id" binding:"required"`
	ScanAt     time.Time `json:"scan_at" binding:"required"`
	IsActive   bool      `json:"is_active" binding:"required"`
}

// VisitHistoryReqEdit 更新访问记录请求
type VisitHistoryReqEdit struct {
	HistoryID  int       `json:"history_id" binding:"required"`
	UserID     int       `json:"user_id"`
	FacilityID int       `json:"facility_id"`
	ScanAt     time.Time `json:"scan_at"`
	IsActive   bool      `json:"is_active"`
}

// VisitHistoryReqList 访问记录分页与搜索请求
type VisitHistoryReqList struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Keyword  string `json:"keyword"`
}

// VisitHistoryDetailRequest 获取单个访问记录请求
type VisitHistoryDetailRequest struct {
	HistoryID int `json:"history_id" binding:"required"`
}
