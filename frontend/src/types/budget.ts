export interface BudgetItem {
  id: number
  project_id: number
  category: string
  budget_amount: number
  actual_amount: number
  variance: number
  remark: string
  created_at: string
  updated_at: string
}

export interface CreateBudgetRequest {
  project_id: number
  category: string
  budget_amount?: number
  actual_amount?: number
  remark?: string
}
