export interface BudgetItem {
  id: number
  project_id: number
  category: string
  budget_amount: number
  /** 手工填写金额：材料类别下为手工补充金额，其他类别即手工实际花费 */
  actual_amount: number
  /** 材料自动花费：已到货/已安装材料的总价，仅 Material 类别非零 */
  material_auto_amount: number
  /** 实际花费合计（材料自动花费 + 手工金额），超支判断以此为准 */
  total_actual_amount: number
  variance: number
  remark: string
  created_at: string
  updated_at: string
}

export interface CreateBudgetRequest {
  project_id: number
  category: string
  budget_amount?: number
  /** 手工填写的实际花费（材料类别下为补充金额） */
  actual_amount?: number
  remark?: string
}
