package utils

import "math"

// BudgetVariance 计算预算差异（实际花费 - 预算金额），保留两位小数。
func BudgetVariance(budgetAmount, actualAmount float64) float64 {
	return Round2(actualAmount - budgetAmount)
}

// MaterialTotal 计算材料总价（数量 x 单价），保留两位小数。
func MaterialTotal(quantity float64, unitPrice float64) float64 {
	return Round2(quantity * unitPrice)
}

// Round2 四舍五入保留两位小数。
func Round2(value float64) float64 {
	return math.Round(value*100) / 100
}

// IsOverBudget 判断是否超支。
func IsOverBudget(variance float64) bool {
	return variance > 0
}
