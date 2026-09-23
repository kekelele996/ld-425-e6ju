package service

import "github.com/home-renovation/platform/internal/model"

// BudgetView 预算项展示视图。
// 在持久化的预算项之上附加材料自动花费，并以实际花费合计作为差异与超支判断依据。
type BudgetView struct {
	model.BudgetItem

	// ManualActualAmount 手工填写的材料实际花费（补充金额）。
	ManualActualAmount float64 `json:"manual_actual_amount"`
	// AutoActualAmount 材料自动花费：项目下已到货/已安装材料的总价汇总。
	AutoActualAmount float64 `json:"auto_actual_amount"`
}

// TotalActualAmount 实际花费合计 = 手工补充 + 材料自动。
func (v *BudgetView) TotalActualAmount() float64 {
	return v.ManualActualAmount + v.AutoActualAmount
}
