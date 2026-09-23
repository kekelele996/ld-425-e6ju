export interface BudgetItem {
  id: number
  project_id: number
  category: string
  budget_amount: number
  /** 实际花费合计：手工补充金额 + 材料自动花费（超支判断以此为准） */
  actual_amount: number
  /** 手工填写的材料实际花费（补充金额） */
  manual_actual_amount: number
  /** 材料自动花费：项目下已到货/已安装材料的总价汇总 */
  auto_actual_amount: number
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
