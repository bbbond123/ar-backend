package model

import (
	"time"
)

type File struct {
	FileID    uint      `gorm:"column:file_id;primaryKey" json:"file_id"`
	FileName  string    `gorm:"column:file_name" json:"file_name"`
	FileType  string    `gorm:"column:file_type" json:"file_type"`
	FileSize  int       `gorm:"column:file_size" json:"file_size"`
	FileData  []byte    `gorm:"column:file_data" json:"-"` // 文件数据一般不直接返回，忽略JSON输出
	Location  string    `gorm:"column:location" json:"location"`
	RelatedID uint      `gorm:"column:related_id" json:"related_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
