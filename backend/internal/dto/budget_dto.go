package dto

import "time"

// CreateBudgetRequest 创建预算项请求。
type CreateBudgetRequest struct {
	ProjectID uint `json:"project_id" binding:"required"`
	Category  string `json:"category" binding:"required,max=32"`
	// BudgetAmount 预算金额。
	BudgetAmount float64 `json:"budget_amount" binding:"gte=0"`
	// ActualAmount 手工填写的实际花费（材料类别下为补充金额，叠加在材料自动花费之上）。
	ActualAmount float64 `json:"actual_amount" binding:"gte=0"`
	Remark       string  `json:"remark" binding:"max=255"`
}

// UpdateBudgetRequest 更新预算项请求。
type UpdateBudgetRequest struct {
	Category     *string  `json:"category" binding:"omitempty,max=32"`
	BudgetAmount *float64 `json:"budget_amount" binding:"omitempty,gte=0"`
	// ActualAmount 手工填写的实际花费（材料类别下为补充金额）。
	ActualAmount *float64 `json:"actual_amount" binding:"omitempty,gte=0"`
	Remark       *string  `json:"remark" binding:"omitempty,max=255"`
}

// BudgetDTO 预算项展示结构。
type BudgetDTO struct {
	ID           uint      `json:"id"`
	ProjectID    uint      `json:"project_id"`
	Category     string    `json:"category"`
	BudgetAmount float64   `json:"budget_amount"`
	// ActualAmount 手工填写金额：材料类别下为手工补充金额，其他类别即手工实际花费。
	ActualAmount float64 `json:"actual_amount"`
	// MaterialAutoAmount 材料自动花费：已到货/已安装材料的总价，仅材料类别非零。
	MaterialAutoAmount float64 `json:"material_auto_amount"`
	// TotalActualAmount 实际花费合计（自动 + 手工），超支判断以此为准。
	TotalActualAmount float64   `json:"total_actual_amount"`
	Variance          float64   `json:"variance"`
	Remark            string    `json:"remark"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
