package model

import "time"

// DesignPhase 设计阶段。
type DesignPhase struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ProjectID    uint      `gorm:"not null;index" json:"project_id"`
	Name         string    `gorm:"size:64;not null" json:"name"`
	DesignerID   uint      `gorm:"index" json:"designer_id"`
	Status       string    `gorm:"size:32;not null;index" json:"status"`
	Version      int       `gorm:"not null;default:1" json:"version"`
	Description  string    `gorm:"type:text" json:"description"`
	FileURLs     string    `gorm:"type:json" json:"file_urls"`
	ReviewComment string   `gorm:"type:text" json:"review_comment"`
	ReviewerID   uint      `json:"reviewer_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名。
func (DesignPhase) TableName() string { return "design_phases" }
