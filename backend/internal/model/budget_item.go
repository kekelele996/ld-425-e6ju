package model

import "time"

// BudgetItem 预算项。
type BudgetItem struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProjectID     uint      `gorm:"not null;index" json:"project_id"`
	Category      string    `gorm:"size:32;not null;index" json:"category"`
	BudgetAmount  float64   `gorm:"type:decimal(14,2);not null;default:0" json:"budget_amount"`
	ActualAmount  float64   `gorm:"type:decimal(14,2);not null;default:0" json:"actual_amount"`
	Variance      float64   `gorm:"type:decimal(14,2);not null;default:0" json:"variance"`
	Remark        string    `gorm:"size:255" json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 指定表名。
func (BudgetItem) TableName() string { return "budget_items" }
