import { create } from 'zustand'
import type { BudgetItem } from '@/types'
import { listBudgets } from '@/api/budget'

interface BudgetState {
  budgets: BudgetItem[]
  loading: boolean
  fetchBudgets: (projectId?: number) => Promise<void>
}

export const useBudgetStore = create<BudgetState>((set) => ({
  budgets: [],
  loading: false,
  fetchBudgets: async (projectId) => {
    set({ loading: true })
    try {
      const result = await listBudgets(projectId ? { project_id: projectId } : { page: 1, page_size: 100 })
      if (Array.isArray(result)) {
        set({ budgets: result, loading: false })
      } else {
        set({ budgets: result.list, loading: false })
      }
    } catch {
      set({ loading: false })
    }
  },
}))
