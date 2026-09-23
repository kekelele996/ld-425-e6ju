package model

import "time"

// MaterialItem 材料清单项。
type MaterialItem struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProjectID      uint      `gorm:"not null;index" json:"project_id"`
	Name           string    `gorm:"size:120;not null" json:"name"`
	Category       string    `gorm:"size:32;not null;index" json:"category"`
	Spec           string    `gorm:"size:120" json:"spec"`
	Brand          string    `gorm:"size:80" json:"brand"`
	Quantity       float64   `gorm:"type:decimal(12,2);not null;default:0" json:"quantity"`
	Unit           string    `gorm:"size:16;not null" json:"unit"`
	UnitPrice      float64   `gorm:"type:decimal(14,2);not null;default:0" json:"unit_price"`
	TotalPrice     float64   `gorm:"type:decimal(14,2);not null;default:0" json:"total_price"`
	PurchaseStatus string    `gorm:"size:32;not null;index" json:"purchase_status"`
	Supplier       string    `gorm:"size:120" json:"supplier"`
	Space          string    `gorm:"size:32;index" json:"space"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 指定表名。
func (MaterialItem) TableName() string { return "material_items" }
