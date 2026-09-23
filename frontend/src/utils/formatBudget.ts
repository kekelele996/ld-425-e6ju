// 预算金额格式化和差异计算。

export function formatCurrency(value: number | undefined | null): string {
  const num = Number(value ?? 0)
  return num.toLocaleString('zh-CN', {
    style: 'currency',
    currency: 'CNY',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}

export function calcVariance(budget: number, actual: number): number {
  return Math.round((actual - budget) * 100) / 100
}

export function isOverBudget(variance: number): boolean {
  return variance > 0
}

export function budgetProgress(budget: number, actual: number): number {
  if (budget <= 0) return 0
  return Math.min(100, Math.round((actual / budget) * 100))
}
