package dto

import "time"

// CreateProjectRequest 创建项目请求。
type CreateProjectRequest struct {
	Name            string  `json:"name" binding:"required,max=120"`
	HouseType       string  `json:"house_type" binding:"required,max=32"`
	Area            float64 `json:"area" binding:"required,gt=0"`
	DecorStyle      string  `json:"decor_style" binding:"required,max=32"`
	Address         string  `json:"address" binding:"max=255"`
	OwnerID         uint    `json:"owner_id" binding:"required"`
	DesignerID      uint    `json:"designer_id"`
	ForemanID       uint    `json:"foreman_id"`
	Status          string  `json:"status" binding:"max=32"`
	ContractAmount  float64 `json:"contract_amount"`
	StartDate       *string `json:"start_date"`
	ExpectedEndDate *string `json:"expected_end_date"`
}

// UpdateProjectRequest 更新项目请求。
type UpdateProjectRequest struct {
	Name            *string  `json:"name" binding:"omitempty,max=120"`
	HouseType       *string  `json:"house_type" binding:"omitempty,max=32"`
	Area            *float64 `json:"area" binding:"omitempty,gt=0"`
	DecorStyle      *string  `json:"decor_style" binding:"omitempty,max=32"`
	Address         *string  `json:"address" binding:"omitempty,max=255"`
	OwnerID         *uint    `json:"owner_id"`
	DesignerID      *uint    `json:"designer_id"`
	ForemanID       *uint    `json:"foreman_id"`
	ContractAmount  *float64 `json:"contract_amount"`
	StartDate       *string  `json:"start_date"`
	ExpectedEndDate *string  `json:"expected_end_date"`
}

// UpdateProjectStatusRequest 项目状态流转请求。
type UpdateProjectStatusRequest struct {
	Status string `json:"status" binding:"required,max=32"`
}

// ProjectDTO 项目展示结构。
type ProjectDTO struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	HouseType       string     `json:"house_type"`
	Area            float64    `json:"area"`
	DecorStyle      string     `json:"decor_style"`
	Address         string     `json:"address"`
	OwnerID         uint       `json:"owner_id"`
	DesignerID      uint       `json:"designer_id"`
	ForemanID       uint       `json:"foreman_id"`
	Status          string     `json:"status"`
	ContractAmount  float64    `json:"contract_amount"`
	StartDate       *time.Time `json:"start_date"`
	ExpectedEndDate *time.Time `json:"expected_end_date"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
