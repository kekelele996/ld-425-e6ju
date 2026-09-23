import { apiDelete, apiGet, apiPost, apiPut } from '@/utils/request'
import type { BudgetItem, CreateBudgetRequest, PageResult } from '@/types'
import { API_PATHS } from '@/constants/apiPaths'

export function listBudgets(params?: Record<string, unknown>): Promise<PageResult<BudgetItem> | BudgetItem[]> {
  return apiGet<PageResult<BudgetItem> | BudgetItem[]>(API_PATHS.budgets, params)
}

export function getBudget(id: number): Promise<BudgetItem> {
  return apiGet<BudgetItem>(`${API_PATHS.budgets}/${id}`)
}

export function createBudget(body: CreateBudgetRequest): Promise<BudgetItem> {
  return apiPost<BudgetItem>(API_PATHS.budgets, body)
}

export function updateBudget(id: number, body: Partial<CreateBudgetRequest>): Promise<BudgetItem> {
  return apiPut<BudgetItem>(`${API_PATHS.budgets}/${id}`, body)
}

export function deleteBudget(id: number): Promise<void> {
  return apiDelete<void>(`${API_PATHS.budgets}/${id}`)
}
