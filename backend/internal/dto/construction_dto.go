package dto

import "time"

// CreateConstructionRequest 创建施工节点请求。
type CreateConstructionRequest struct {
	ProjectID        uint     `json:"project_id" binding:"required"`
	Name             string   `json:"name" binding:"required,max=64"`
	PlannedStartDate *string  `json:"planned_start_date"`
	PlannedEndDate   *string  `json:"planned_end_date"`
}

// UpdateConstructionRequest 更新施工节点请求。
type UpdateConstructionRequest struct {
	Name             *string `json:"name" binding:"omitempty,max=64"`
	PlannedStartDate *string `json:"planned_start_date"`
	PlannedEndDate   *string `json:"planned_end_date"`
}

// UpdateConstructionStatusRequest 施工状态流转请求。
type UpdateConstructionStatusRequest struct {
	Status string `json:"status" binding:"required,max=32"`
}

// AcceptConstructionRequest 施工验收请求。
type AcceptConstructionRequest struct {
	Accepted bool     `json:"accepted"`
	Photos   []string `json:"photos"`
	Note     string   `json:"note" binding:"max=500"`
}

// ConstructionDTO 施工节点展示结构。
type ConstructionDTO struct {
	ID               uint       `json:"id"`
	ProjectID        uint       `json:"project_id"`
	Name             string     `json:"name"`
	PlannedStartDate *time.Time `json:"planned_start_date"`
	PlannedEndDate   *time.Time `json:"planned_end_date"`
	ActualStartDate  *time.Time `json:"actual_start_date"`
	ActualEndDate    *time.Time `json:"actual_end_date"`
	Status           string     `json:"status"`
	AcceptanceStatus string     `json:"acceptance_status"`
	AcceptancePhotos []string   `json:"acceptance_photos"`
	AcceptanceNote   string     `json:"acceptance_note"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
