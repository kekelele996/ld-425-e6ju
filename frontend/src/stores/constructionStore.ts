import { create } from 'zustand'
import type { ConstructionNode } from '@/types'
import { listConstructions } from '@/api/construction'

interface ConstructionState {
  nodes: ConstructionNode[]
  loading: boolean
  fetchNodes: (projectId?: number) => Promise<void>
}

export const useConstructionStore = create<ConstructionState>((set) => ({
  nodes: [],
  loading: false,
  fetchNodes: async (projectId) => {
    set({ loading: true })
    try {
      const result = await listConstructions(projectId ? { project_id: projectId } : { page: 1, page_size: 100 })
      if (Array.isArray(result)) {
        set({ nodes: result, loading: false })
      } else {
        set({ nodes: result.list, loading: false })
      }
    } catch {
      set({ loading: false })
    }
  },
}))
