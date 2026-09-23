package dto

import "time"

// CreateBudgetRequest 创建预算项请求。
type CreateBudgetRequest struct {
	ProjectID    uint    `json:"project_id" binding:"required"`
	Category     string  `json:"category" binding:"required,max=32"`
	BudgetAmount float64 `json:"budget_amount" binding:"gte=0"`
	ActualAmount float64 `json:"actual_amount" binding:"gte=0"`
	Remark       string  `json:"remark" binding:"max=255"`
}

// UpdateBudgetRequest 更新预算项请求。
type UpdateBudgetRequest struct {
	Category     *string  `json:"category" binding:"omitempty,max=32"`
	BudgetAmount *float64 `json:"budget_amount" binding:"omitempty,gte=0"`
	ActualAmount *float64 `json:"actual_amount" binding:"omitempty,gte=0"`
	Remark       *string  `json:"remark" binding:"omitempty,max=255"`
}

// BudgetDTO 预算项展示结构。
type BudgetDTO struct {
	ID           uint      `json:"id"`
	ProjectID    uint      `json:"project_id"`
	Category     string    `json:"category"`
	BudgetAmount float64   `json:"budget_amount"`
	ActualAmount float64   `json:"actual_amount"`
	Variance     float64   `json:"variance"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
