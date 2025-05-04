package model

import "time"

// Menu 表示数据库中的 menus 表
type Menu struct {
	MenuID       int        `gorm:"column:menu_id;primaryKey" json:"menu_id"`
	MenuName     string     `gorm:"column:menu_name;type:varchar(100);not null" json:"menu_name"`
	MenuCode     string     `gorm:"column:menu_code;type:varchar(50);not null" json:"menu_code"`
	DisplayOrder *int       `gorm:"column:display_order" json:"display_order"`
	IsActive     bool       `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// MenuReqCreate 新建菜单请求
type MenuReqCreate struct {
	MenuName     string `json:"menu_name" binding:"required"`
	MenuCode     string `json:"menu_code" binding:"required"`
	DisplayOrder *int   `json:"display_order"`
	IsActive     bool   `json:"is_active" binding:"required"`
}

// MenuReqEdit 更新菜单请求
type MenuReqEdit struct {
	MenuID       int    `json:"menu_id" binding:"required"`
	MenuName     string `json:"menu_name"`
	MenuCode     string `json:"menu_code"`
	DisplayOrder *int   `json:"display_order"`
	IsActive     bool   `json:"is_active"`
}

// MenuReqList 菜单分页与搜索请求
type MenuReqList struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Keyword  string `json:"keyword"`
}

// MenuDetailRequest 获取单个菜单请求
type MenuDetailRequest struct {
	MenuID int `json:"menu_id" binding:"required"`
}
