package model

import (
	"time"
)

// Facility 表示数据库中的 facilities 表
// 数据库表结构
// CREATE TABLE facilities (
//
//	facility_id SERIAL PRIMARY KEY,
//	facility_name VARCHAR(255) NOT NULL,
//	location VARCHAR(255) NOT NULL,
//	description_text TEXT,
//	latitude DECIMAL(10,6) NOT NULL,
//	longitude DECIMAL(10,6) NOT NULL,
//	person_id INTEGER,
//	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
//	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
//
// );
type Facility struct {
	FacilityID      int       `gorm:"column:facility_id;primaryKey" json:"facility_id"`                     // 设施ID
	FacilityName    string    `gorm:"column:facility_name;type:varchar(255);not null" json:"facility_name"` // 设施名
	Location        string    `gorm:"column:location;type:varchar(255);not null" json:"location"`           // 所在地
	DescriptionText string    `gorm:"column:description_text;type:text" json:"description"`                 // 设施描述
	Latitude        float64   `gorm:"column:latitude;type:decimal(10,6);not null" json:"latitude"`          // 纬度
	Longitude       float64   `gorm:"column:longitude;type:decimal(10,6);not null" json:"longitude"`        // 经度
	PersonID        *int      `gorm:"column:person_id" json:"person_id"`                                    // 相关人物ID（可选）
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// FacilityReqCreate 用于创建设施时的请求参数
type FacilityReqCreate struct {
	FacilityName    string  `json:"facility_name" binding:"required"` // 设施名
	Location        string  `json:"location" binding:"required"`      // 所在地
	DescriptionText string  `json:"description"`                      // 描述
	Latitude        float64 `json:"latitude" binding:"required"`      // 纬度
	Longitude       float64 `json:"longitude" binding:"required"`     // 经度
	PersonID        *int    `json:"person_id"`                        // 相关人物ID（可选）
}

// FacilityReqEdit 用于编辑设施时的请求参数
type FacilityReqEdit struct {
	FacilityID      int     `json:"facility_id" binding:"required"`   // 设施ID
	FacilityName    string  `json:"facility_name" binding:"required"` // 设施名
	Location        string  `json:"location" binding:"required"`      // 所在地
	DescriptionText string  `json:"description"`                      // 描述
	Latitude        float64 `json:"latitude" binding:"required"`      // 纬度
	Longitude       float64 `json:"longitude" binding:"required"`     // 经度
	PersonID        *int    `json:"person_id"`                        // 相关人物ID（可选）
}

// FacilityReqList 用于分页与查询设施时的请求参数
type FacilityReqList struct {
	Page     int    `json:"page" binding:"required"`      // 页码
	PageSize int    `json:"page_size" binding:"required"` // 每页数量
	Keyword  string `json:"keyword"`                      // 关键字（设施名模糊搜索）
}

// FacilityCreateRequest 兼容风格，新建请求（备用，与 ReqCreate 作用相同）
type FacilityCreateRequest = FacilityReqCreate

// FacilityUpdateRequest 兼容风格，更新请求（备用，与 ReqEdit 作用相同）
type FacilityUpdateRequest = FacilityReqEdit

// FacilityQueryRequest 兼容风格，分页查询请求（备用，与 ReqList 作用相同）
type FacilityQueryRequest = FacilityReqList

// FacilityDetailRequest 获取单个设施请求
type FacilityDetailRequest struct {
	FacilityID int `json:"facility_id" binding:"required"` // 设施ID
}
