package model

import "time"

// File 表示数据库中的 files 表
type File struct {
	FileID    int       `gorm:"column:file_id;primaryKey" json:"file_id"`
	FileName  string    `gorm:"column:file_name;type:varchar(255);not null" json:"file_name"`
	FileType  string    `gorm:"column:file_type;type:varchar(50);not null" json:"file_type"`
	FileSize  int       `gorm:"column:file_size" json:"file_size"`
	FileData  []byte    `gorm:"column:file_data;not null" json:"file_data"`
	Location  string    `gorm:"column:location;type:varchar(255);not null" json:"location"`
	RelatedID int       `gorm:"column:related_id;not null" json:"related_id"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// FileReqCreate 新建文件请求
type FileReqCreate struct {
	FileName  string `json:"file_name" binding:"required"`
	FileType  string `json:"file_type" binding:"required"`
	FileSize  int    `json:"file_size"`
	FileData  []byte `json:"file_data" binding:"required"`
	Location  string `json:"location" binding:"required"`
	RelatedID int    `json:"related_id" binding:"required"`
}

// FileReqEdit 更新文件请求
type FileReqEdit struct {
	FileID    int    `json:"file_id" binding:"required"`
	FileName  string `json:"file_name"`
	FileType  string `json:"file_type"`
	FileSize  int    `json:"file_size"`
	FileData  []byte `json:"file_data"`
	Location  string `json:"location"`
	RelatedID int    `json:"related_id"`
}

// FileReqList 文件分页与搜索请求
type FileReqList struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Keyword  string `json:"keyword"`
}

// FileDetailRequest 获取单个文件请求
type FileDetailRequest struct {
	FileID int `json:"file_id" binding:"required"`
}
