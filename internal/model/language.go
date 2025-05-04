package model

import (
	"time"
)

// Facility 表示数据库中的 facilities 表
type Language struct {
	FacilityID      int       `gorm:"column:facility_id;primaryKey"`                        // 施設ID: 主键
	FacilityName    string    `gorm:"column:facility_name;type:varchar(255);not null"`      // 施設名（全角）
	Location        string    `gorm:"column:location;type:varchar(255);not null"`           // 所在地（全角）
	DescriptionText string    `gorm:"column:description_text;type:text"`                    // 説明（全角）
	Latitude        float64   `gorm:"column:latitude;type:decimal(10,6);not null"`          // 緯度
	Longitude       float64   `gorm:"column:longitude;type:decimal(10,6);not null"`         // 経度
	PersonID        *int      `gorm:"column:person_id"`                                     // 関連人物ID（允许为NULL）
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"` // 作成日
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"` // 更新日
}
