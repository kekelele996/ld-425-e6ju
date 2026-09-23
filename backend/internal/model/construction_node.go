package model

import "time"

// ConstructionNode 施工节点。
type ConstructionNode struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	ProjectID         uint       `gorm:"not null;index" json:"project_id"`
	Name              string     `gorm:"size:64;not null" json:"name"`
	PlannedStartDate  *time.Time `gorm:"type:date" json:"planned_start_date"`
	PlannedEndDate    *time.Time `gorm:"type:date" json:"planned_end_date"`
	ActualStartDate   *time.Time `gorm:"type:date" json:"actual_start_date"`
	ActualEndDate     *time.Time `gorm:"type:date" json:"actual_end_date"`
	Status            string     `gorm:"size:32;not null;index" json:"status"`
	AcceptanceStatus  string     `gorm:"size:32;not null;index" json:"acceptance_status"`
	AcceptancePhotos  string     `gorm:"type:json" json:"acceptance_photos"`
	AcceptanceNote    string     `gorm:"type:text" json:"acceptance_note"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// TableName 指定表名。
func (ConstructionNode) TableName() string { return "construction_nodes" }
