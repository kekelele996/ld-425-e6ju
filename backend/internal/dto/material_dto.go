package dto

import "time"

// CreateMaterialRequest 创建材料请求。
type CreateMaterialRequest struct {
	ProjectID uint    `json:"project_id" binding:"required"`
	Name      string  `json:"name" binding:"required,max=120"`
	Category  string  `json:"category" binding:"required,max=32"`
	Spec      string  `json:"spec" binding:"max=120"`
	Brand     string  `json:"brand" binding:"max=80"`
	Quantity  float64 `json:"quantity" binding:"required,gte=0"`
	Unit      string  `json:"unit" binding:"required,max=16"`
	UnitPrice float64 `json:"unit_price" binding:"gte=0"`
	Supplier  string  `json:"supplier" binding:"max=120"`
	Space     string  `json:"space" binding:"max=32"`
}

// UpdateMaterialRequest 更新材料请求。
type UpdateMaterialRequest struct {
	Name      *string  `json:"name" binding:"omitempty,max=120"`
	Category  *string  `json:"category" binding:"omitempty,max=32"`
	Spec      *string  `json:"spec" binding:"omitempty,max=120"`
	Brand     *string  `json:"brand" binding:"omitempty,max=80"`
	Quantity  *float64 `json:"quantity" binding:"omitempty,gte=0"`
	Unit      *string  `json:"unit" binding:"omitempty,max=16"`
	UnitPrice *float64 `json:"unit_price" binding:"omitempty,gte=0"`
	Supplier  *string  `json:"supplier" binding:"omitempty,max=120"`
	Space     *string  `json:"space" binding:"omitempty,max=32"`
}

// UpdateMaterialStatusRequest 采购状态流转请求。
type UpdateMaterialStatusRequest struct {
	Status string `json:"status" binding:"required,max=32"`
}

// MaterialDTO 材料展示结构。
type MaterialDTO struct {
	ID             uint      `json:"id"`
	ProjectID      uint      `json:"project_id"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Spec           string    `json:"spec"`
	Brand          string    `json:"brand"`
	Quantity       float64   `json:"quantity"`
	Unit           string    `json:"unit"`
	UnitPrice      float64   `json:"unit_price"`
	TotalPrice     float64   `json:"total_price"`
	PurchaseStatus string    `json:"purchase_status"`
	Supplier       string    `json:"supplier"`
	Space          string    `json:"space"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
