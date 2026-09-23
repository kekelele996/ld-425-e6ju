import { create } from 'zustand'
import type { DesignPhase } from '@/types'
import { listDesigns } from '@/api/design'

interface DesignState {
  designs: DesignPhase[]
  loading: boolean
  fetchDesigns: (projectId?: number) => Promise<void>
}

export const useDesignStore = create<DesignState>((set) => ({
  designs: [],
  loading: false,
  fetchDesigns: async (projectId) => {
    set({ loading: true })
    try {
      const result = await listDesigns(projectId ? { project_id: projectId } : { page: 1, page_size: 100 })
      if (Array.isArray(result)) {
        set({ designs: result, loading: false })
      } else {
        set({ designs: result.list, loading: false })
      }
    } catch {
      set({ loading: false })
    }
  },
}))
