import { create } from 'zustand'
import type { RenovationProject } from '@/types'
import { listProjects } from '@/api/project'

interface ProjectState {
  projects: RenovationProject[]
  total: number
  loading: boolean
  fetchProjects: () => Promise<void>
}

export const useProjectStore = create<ProjectState>((set) => ({
  projects: [],
  total: 0,
  loading: false,
  fetchProjects: async () => {
    set({ loading: true })
    try {
      const result = await listProjects({ page: 1, page_size: 50 })
      if (Array.isArray(result)) {
        set({ projects: result, total: result.length, loading: false })
      } else {
        set({ projects: result.list, total: result.total, loading: false })
      }
    } catch {
      set({ loading: false })
    }
  },
}))
